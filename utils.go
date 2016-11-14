package web

import (
	"net"
	"time"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
	keepAliveDuration time.Duration
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(ln.keepAliveDuration)
	return tc, nil
}
