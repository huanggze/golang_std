package net

import (
	"runtime"
	"syscall"
)

// Boolean to int.
func boolint(b bool) int {
	if b {
		return 1
	}
	return 0
}

func setKeepAlive(fd *netFD, keepalive bool) error {
	err := fd.pfd.SetsockoptInt(syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, boolint(keepalive))
	runtime.KeepAlive(fd)
	return wrapSyscallError("setsockopt", err)
}
