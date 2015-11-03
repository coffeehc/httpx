// websocket
package web

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

func checkOrigin(config *websocket.Config, req *http.Request) (err error) {
	config.Origin, err = websocket.Origin(config, req)
	if err == nil && config.Origin == nil {
		return fmt.Errorf("null origin")
	}
	return err
}

type webScoketAdapter struct {
	handler WebScoketHandler
}

func (this *webScoketAdapter) webScoketHandlerAdapter(request *http.Request, pathFragments map[string]string, reply *Reply) {
	reply.startWebSocket()
	server := websocket.Server{Handler: func(conn *websocket.Conn) {
		defer func() {
			conn.Close()
		}()
		this.handler(request, pathFragments, &WebSocketReply{conn})
	}, Handshake: checkOrigin}
	server.ServeHTTP(reply.w, request)
}

type WebScoketHandler func(request *http.Request, pathFragments map[string]string, reply *WebSocketReply)

type WebSocketReply struct {
	conn *websocket.Conn
}

func (this *WebSocketReply) Read(msg []byte) (n int, err error) {
	return this.conn.Read(msg)
}

func (this *WebSocketReply) Write(msg []byte) (n int, err error) {
	return this.conn.Write(msg)
}

func (this *WebSocketReply) Close() error {
	return this.conn.Close()
}

func (this *WebSocketReply) IsClientConn() bool {
	return this.conn.IsClientConn()
}
func (this *WebSocketReply) IsServerConn() bool {
	return this.conn.IsServerConn()
}

func (this *WebSocketReply) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *WebSocketReply) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}
func (this *WebSocketReply) SetDeadline(t time.Time) error {
	return this.conn.SetDeadline(t)
}
func (this *WebSocketReply) SetReadDeadline(t time.Time) error {
	return this.conn.SetReadDeadline(t)
}
func (this *WebSocketReply) SetWriteDeadline(t time.Time) error {
	return this.conn.SetWriteDeadline(t)
}
func (this *WebSocketReply) Config() *websocket.Config {
	return this.conn.Config()
}
