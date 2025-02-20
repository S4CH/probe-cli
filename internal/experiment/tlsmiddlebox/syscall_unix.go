//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || nacl || netbsd || openbsd || solaris

package tlsmiddlebox

//
// syscall utilities for dialerTTLWrapperConn
//

import (
	"net"
	"syscall"
)

// SetTTL sets the IP TTL field for the underlying net.TCPConn
func (c *dialerTTLWrapperConn) SetTTL(ttl int) error {
	conn := c.Conn
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return errInvalidConnWrapper
	}
	rawConn, err := tcpConn.SyscallConn()
	if err != nil {
		return err
	}
	rawErr := rawConn.Control(func(fd uintptr) {
		err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_TTL, ttl)
	})
	// The syscall err is given a higher priority and returned early if non-nil
	if err != nil {
		return err
	}
	return rawErr
}

// GetSoErr fetches the SO_ERROR value to look for soft ICMP errors in TCP
func (c *dialerTTLWrapperConn) GetSoErr() (errno int, err error) {
	conn := c.Conn
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return 0, errInvalidConnWrapper
	}
	rawConn, err := tcpConn.SyscallConn()
	if err != nil {
		return 0, errInvalidConnWrapper
	}
	rawErr := rawConn.Control(func(fd uintptr) {
		errno, err = syscall.GetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_ERROR)
	})
	if rawErr != nil {
		return 0, rawErr
	}
	return
}
