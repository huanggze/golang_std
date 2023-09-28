package net

import (
	"context"
	"syscall"
)

func favoriteAddrFamily(network string, laddr, raddr sockaddr, mode string) (family int, ipv6only bool) {
	return syscall.AF_INET, false
}

func internetSocket(ctx context.Context, net string, laddr, raddr sockaddr, sotype, proto int, mode string) (fd *netFD, err error) {
	family, ipv6only := favoriteAddrFamily(net, laddr, raddr, mode)
	return socket(ctx, net, family, sotype, proto, ipv6only, laddr, raddr)
}
