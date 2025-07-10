package main

import (
	inout "cmdline/InputOutput"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

var commandToSyscall = map[string]string{
	// "mkdir":  "mkdir",
	// "rm":     "unlink",
	"rmdir": "rmdir",
	"touch": "open",
	"ls":    "getdents",
	"tree":  "getdents (recursive)",
	// "cat":    "read",
	// "echo":   "write",
	// "whoami": "getuid",
	"chmod": "chmod",
	"chown": "chown",
	"stat":  "stat",
	"fstat": "fstat",
	// "lstat": "lstat",
	// "mv":    "rename",
}

var commandHandlers = map[string]func([]string) error{
	"cd": func(args []string) error {
		if len(args) < 1 {
			return inout.WriteStdout("cd: Not enough arguments provided\n")
		}
		return syscall.Chdir(args[0])
	},
	"mkdir": func(args []string) error {
		if len(args) < 1 {
			return inout.WriteStdout("mkdir: Not enough arguments provided\n")
		}
		return syscall.Mkdir(args[0], 0755)
	},
	"rm": func(args []string) error {
		if len(args) < 1 {
			return inout.WriteStdout("rm: Not enough arguments provided\n")
		}
		return syscall.Unlink(args[0])
	},
	"echo": func(args []string) error {
		var str string

		for i := range args {
			str += args[i] + " "
		}
		str += "\n"
		return inout.WriteStdout(str)

	},
	"whoami": func(args []string) error {
		uid := syscall.Getuid()
		u, err := user.LookupId(strconv.Itoa(uid))
		if err != nil {
			res := fmt.Sprintf("Error:: %d", uid)
			inout.WriteStdout(res)
		}
		res := fmt.Sprintf(u.Username + "\n")
		return inout.WriteStdout(res)
	},
	"pwd": func(args []string) error {
		buf := make([]byte, 512)
		n, err := syscall.Getcwd(buf)

		if err != nil {
			inout.WriteStdout("pwd:: Error unknown")
		}

		res := fmt.Sprintf(string(buf[:n]) + "\n")
		return inout.WriteStdout(res)
	},

	"cat": func(args []string) error {
		buf := make([]byte, 4096)

		if len(args) < 1 {
			return inout.WriteStdout("cat:: Not enough arguments provided\n")
		}

		for i := range args {

			fd, err := syscall.Open(args[i], syscall.O_RDONLY, 0)

			if err != nil {
				return inout.WriteStdout("Invalid path")
			}

			defer syscall.Close(fd)

			for {

				n, err := syscall.Read(fd, buf)

				if err != nil {
					errr := fmt.Sprintf("cat:: %s", err.Error())

					return inout.WriteStdout(errr)
				}

				if n == 0 {
					break
				}

				err = inout.WriteStdout(string(buf[:n]))

				if err != nil {
					return err
				}
			}

			inout.WriteStdout("\n")
		}

		return nil

	},
	"mv": func(args []string) error {
		if len(args) < 2 {
			return inout.WriteStdout("mv:: Not enough arguments provided\n")
		}

		err := syscall.Rename(args[0], args[1])

		if err != nil {

			errr := fmt.Sprintf("mv::%s\n", err.Error())
			return inout.WriteStdout(errr)
		}

		return nil

	},
	"lstat": func(args []string) error {

		if len(args) < 1 {
			return inout.WriteStdout("lstat:: Not enough arguments provided\n")
		}
		var stat syscall.Stat_t

		err := syscall.Lstat(args[0], &stat)

		if err != nil {
			errr := fmt.Sprintf("lstat::%s\n", err.Error())
			return inout.WriteStdout(errr)
		}

		output := fmt.Sprintf(
			"File: %s\nSize: %d bytes\nUID: %d\nGID: %d\nMode: %o\n",
			args[0],
			stat.Size,
			stat.Uid,
			stat.Gid,
			stat.Mode,
		)

		return inout.WriteStdout(output)
	},
	"getpid": func(args []string) error {
		res := syscall.Getpid()
		return inout.WriteStdout(strconv.Itoa(res) + "\n")
	},
	"getppid": func(args []string) error {
		res := syscall.Getppid()
		return inout.WriteStdout(strconv.Itoa(res) + "\n")

	},
	"kill": func(args []string) error {
		if len(args) < 1 {
			return inout.WriteStdout("kill: Not enough arguments provided\n")
		}

		pid, err := strconv.Atoi(args[0])
		if err != nil {
			return inout.WriteStdout("kill: Invalid PID\n")
		}

		sig := syscall.SIGTERM

		if len(args) > 1 {
			switch args[1] {
			case "-9", "SIGKILL":
				sig = syscall.SIGTERM
			case "-15", "SIGTERM":
				sig = syscall.SIGTERM
			default:
				return inout.WriteStdout("kill: Unknown signal\n")
			}

		}

		err = syscall.Kill(pid, sig)

		if err != nil {
			return inout.WriteStdout(fmt.Sprintf("kill: %s\n", err.Error()))
		}
		return nil
	},

	"getIP": func(args []string) error {
		fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)
		if err != nil {
			return inout.WriteStdout(fmt.Sprintf("ip::%s\n", err.Error()))
		}

		defer syscall.Close(fd)

		var addr syscall.SockaddrInet4
		addr.Port = 80
		addr.Addr = [4]byte{8, 8, 8, 8}

		err = syscall.Connect(fd, &addr)
		if err != nil {
			return inout.WriteStdout(fmt.Sprintf("getsocketname: %s\n", err.Error()))
		}

		sa, err := syscall.Getsockname(fd)

		if err != nil {
			return inout.WriteStdout(fmt.Sprintf("getsocketname: %s\n", err.Error()))
		}

		switch addr := sa.(type) {
		case *syscall.SockaddrInet4:
			fmt.Printf("IPv4: %v\n", addr.Addr)
		case *syscall.SockaddrInet6:
			fmt.Printf("IPv6: %v\n", addr.Addr)
		case *syscall.SockaddrUnix:
			fmt.Printf("Unix socket path: %s\n", addr.Name)
		default:
			fmt.Printf("Unknown socket type\n")
		}

		return nil
	},
	"time": func(args []string) error {
		var t syscall.Time_t
		t, err := syscall.Time(&t)

		if err != nil {
			return inout.WriteStdout(fmt.Sprintf("time::%s\n", err.Error()))
		}

		res := fmt.Sprintf("%s \n", strconv.FormatInt(int64(t), 10))
		return inout.WriteStdout(res)
	},
	"gettimeofday": func(args []string) error {
		var tod syscall.Timeval
		err := syscall.Gettimeofday(&tod)

		if err != nil {
			return inout.WriteStdout(fmt.Sprintf("time::%s\n", err.Error()))
		}

		return inout.WriteStdout(fmt.Sprintf("%d\n", tod))

	},
	// "fork": func(args []string) error {
	// 	if len(args) == 0 {
	// 		return fmt.Errorf("no command provided")
	// 	}

	// 	path, err := exec.LookPath(args[0])

	// 	attr := &syscall.ProcAttr{
	// 		Files: []uintptr{0, 1, 2},
	// 	}

	// 	if err != nil {
	// 		return fmt.Errorf("command not found: %s", args[0])
	// 	}

	// 	_, err = syscall.ForkExec(path, args, attr)
	// 	if err != nil {
	// 		return fmt.Errorf("ForkExec failed: %v", err)
	// 	}

	// 	return nil
	// },
	"ls": func(args []string) error {
		fd, err := syscall.Openat(syscall., ".", syscall.O_RDONLY|syscall.O_NONBLOCK|syscall.O_CLOEXEC, 0)

		if err != nil {
			return fmt.Errorf("ForkExec failed: %v", err)
		}

		return nil
	},
}

func main() {
	args := os.Args

	if len(args) < 2 {
		inout.WriteStdout("Usage: <command> [args...]\n")
		return
	}

	cmd := args[1]
	cmdArgs := args[2:]

	handler, ok := commandHandlers[cmd]
	if !ok {
		inout.WriteStdout("Unknown command: " + cmd + "\n")
		return
	}

	err := handler(cmdArgs)
	if err != nil {
		inout.WriteStdout("Error: " + err.Error() + "\n")
	}
}

func init() {
	commandHandlers["netstat"] = func(args []string) error {
		catFunc, ok := commandHandlers["cat"]
		if !ok {
			return inout.WriteStdout("netstat:: 'cat' command not found\n")
		}
		return catFunc([]string{"/proc/net/tcp"})
	}
}
