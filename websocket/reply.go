package websocket

import (
	"golang.org/x/net/websocket"
	"net"
	"time"
)

type WebSocket struct {
	conn *websocket.Conn
}

func (this *WebSocket) Read(msg []byte) (n int, err error) {
	return this.conn.Read(msg)
}

func (this *WebSocket) Write(msg []byte) (n int, err error) {
	return this.conn.Write(msg)
}

func (this *WebSocket) Close() error {
	return this.conn.Close()
}

func (this *WebSocket) IsClientConn() bool {
	return this.conn.IsClientConn()
}
func (this *WebSocket) IsServerConn() bool {
	return this.conn.IsServerConn()
}

func (this *WebSocket) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *WebSocket) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}
func (this *WebSocket) SetDeadline(t time.Time) error {
	return this.conn.SetDeadline(t)
}
func (this *WebSocket) SetReadDeadline(t time.Time) error {
	return this.conn.SetReadDeadline(t)
}
func (this *WebSocket) SetWriteDeadline(t time.Time) error {
	return this.conn.SetWriteDeadline(t)
}
func (this *WebSocket) Config() *websocket.Config {
	return this.conn.Config()
}
