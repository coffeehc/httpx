// server_test
package web

import (
	"errors"
	"fmt"
	"github.com/coffeehc/logger"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	logger.InitLogger()
	server := NewServer(nil)
	server.Regedit("/a/{name}/123", GET, Service)
	server.Regedit("/a/123/{name}", GET, testService)
	server.RegeditWebSocket("/api/websocket", WebsocketTest)
	server.AddLastFilter("/*", AccessLogFilter)
	server.Start()
	time.Sleep(time.Second * 160)
	server.Stop()
}

func WebsocketTest(request *http.Request, param map[string]string, reply *WebSocketReply) {
	for i := 0; i < 1000; i++ {
		fmt.Fprint(reply, "hello %d", i)
	}
	reply.Close()

}

func Service(request *http.Request, param map[string]string, reply *Reply) {
	reply.With("123" + param["name"])
	panic(errors.New("test error"))
}

func testService(request *http.Request, param map[string]string, reply *Reply) {
	reply.With(param["name"])
}
