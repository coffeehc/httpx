// websocket
package websocket

import (
	"fmt"
	"net/http"

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
		this.handler(request, pathFragments, &WebSocket{conn})
	}, Handshake: checkOrigin}
	server.ServeHTTP(reply.w, request)
}

type WebScoketHandler func(request *http.Request, pathFragments map[string]string, reply *WebSocket)
