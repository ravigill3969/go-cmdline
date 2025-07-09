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
	"cd":     "chdir",
	"mkdir":  "mkdir",
	"rm":     "unlink",
	"rmdir":  "rmdir",
	"touch":  "open",
	"ls":     "getdents",
	"tree":   "getdents (recursive)",
	"cat":    "read",
	"echo":   "write",
	"whoami": "getuid",
	"chmod":  "chmod",
	"chown":  "chown",
	"stat":   "stat",
	"fstat":  "fstat",
	"lstat":  "lstat",
	"mv":     "rename",
	"pwd":    "getcwd",
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
			return inout.WriteStdout("cat::Not enough arguments provided\n")
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
