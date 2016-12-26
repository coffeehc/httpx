package web

import (
	"errors"
	"net"
	"time"
)

type tcpKeepAliveListener struct {
	net.Listener
	keepAliveDuration time.Duration
}

func (ln *tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	conn, err := ln.Accept()
	if err != nil {
		return nil, err
	}
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(ln.keepAliveDuration)
		return conn, nil
	}
	return nil, errors.New("net.Listener only support tcp")
}
