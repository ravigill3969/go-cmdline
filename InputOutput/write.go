package inout

import "syscall"

func WriteStdout(s string) error {
	data := []byte(s)
	n, err := syscall.Write(1, data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return syscall.EIO
	}
	return nil
}
