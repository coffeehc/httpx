package websocket

import (
	"github.com/coffeehc/web"
	"golang.org/x/net/websocket"
)

func RegeditWebSocket(server web.Server, path string, service websocket.Handler) error {
	return server.RegeditHttpHandler(path, web.GET, service)
}
