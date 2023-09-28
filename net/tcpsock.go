package net

import (
	"github.com/huanggze/src/internal/itoa"
	"syscall"
	"time"
)

// TCPAddr represents the address of a TCP end point.
type TCPAddr struct {
	IP   IP
	Port int
	Zone string // IPv6 scoped addressing zone
}

// Network returns the address's network name, "tcp".
func (a *TCPAddr) Network() string { return "tcp" }

func (a *TCPAddr) String() string {
	if a == nil {
		return "<nil>"
	}
	ip := ipEmptyString(a.IP)
	if a.Zone != "" {
		return JoinHostPort(ip+"%"+a.Zone, itoa.Itoa(a.Port))
	}
	return JoinHostPort(ip, itoa.Itoa(a.Port))
}

// TCPConn is an implementation of the Conn interface for TCP network
// connections.
type TCPConn struct {
	conn
}

// SetNoDelay controls whether the operating system should delay
// packet transmission in hopes of sending fewer packets (Nagle's
// algorithm).  The default is true (no delay), meaning that data is
// sent as soon as possible after a Write.
func (c *TCPConn) SetNoDelay(noDelay bool) error {
	if !c.ok() {
		return syscall.EINVAL
	}
	if err := setNoDelay(c.fd, noDelay); err != nil {
		return &OpError{Op: "set", Net: c.fd.net, Source: c.fd.laddr, Addr: c.fd.raddr, Err: err}
	}
	return nil
}

func newTCPConn(fd *netFD) *TCPConn {
	c := &TCPConn{conn{fd}}
	setNoDelay(c.fd, true)
	return c
}

// roundDurationUp rounds d to the next multiple of to.
func roundDurationUp(d time.Duration, to time.Duration) time.Duration {
	return (d + to - 1) / to
}
