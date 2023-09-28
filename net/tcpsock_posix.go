package net

import (
	"context"
	"syscall"
)

func sockaddrToTCP(sa syscall.Sockaddr) Addr {
	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		return &TCPAddr{IP: sa.Addr[0:], Port: sa.Port}
	}
	return nil
}

func (sd *sysDialer) dialTCP(ctx context.Context, laddr, raddr *TCPAddr) (*TCPConn, error) {
	return sd.doDialTCP(ctx, laddr, raddr)
}

func (sd *sysDialer) doDialTCP(ctx context.Context, laddr, raddr *TCPAddr) (*TCPConn, error) {
	fd, err := internetSocket(ctx, sd.network, laddr, raddr, syscall.SOCK_STREAM, 0, "dial")
	if err != nil {
		return nil, err
	}
	return newTCPConn(fd), nil
}
