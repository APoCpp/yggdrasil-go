// +build linux

package yggdrasil

import (
	"syscall"

	"golang.org/x/sys/unix"
)

// WARNING: This context is used both by net.Dialer and net.Listen in tcp.go

func (t *tcp) tcpContext(network, address string, c syscall.RawConn) error {
	var control error
	var bbr error

	control = c.Control(func(fd uintptr) {
		// sys/socket.h: #define	SO_RECV_ANYIF	0x1104
		bbr = unix.SetsockoptString(int(fd), unix.IPPROTO_TCP, unix.TCP_CONGESTION, "bbr")
	})

	// Log any errors
	if bbr != nil {
		t.link.core.log.Debugln("Failed to set tcp_congestion_control to bbr for socket, SetsockoptString error:", bbr)
	}
	if control != nil {
		t.link.core.log.Debugln("Failed to set tcp_congestion_control to bbr for socket, Control error:", control)
	}

	// Return nil because errors here are not considered fatal for the connection, it just means congestion control is suboptimal
	return nil
}
