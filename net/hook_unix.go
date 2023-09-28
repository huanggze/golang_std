//go:build unix || (js && wasm)

package net

import "syscall"

var (
	connectFunc       func(int, syscall.Sockaddr) error = syscall.Connect
	listenFunc        func(int, int) error              = syscall.Listen
	getsockoptIntFunc func(int, int, int) (int, error)  = syscall.GetsockoptInt
)
