package net

import "syscall"

// A sockaddr represents a TCP, UDP, IP or Unix network endpoint
// address that can be converted into a syscall.Sockaddr.
type sockaddr interface {
	// sockaddr returns the address converted into a syscall
	// sockaddr type that implements syscall.Sockaddr
	// interface. It returns a nil interface when the address is
	// nil.
	sockaddr(family int) (syscall.Sockaddr, error)
}
