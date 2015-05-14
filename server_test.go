// server_test
package web

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	server := NewServer(nil)
	server.Regedit("/a/{name}/123", GET, Service)
	server.Regedit("/a/123/{name}", GET, testService)
	server.AddFilter("/*", AccessLogFilter)
	server.Start()
	time.Sleep(time.Second * 20)
	server.Stop()
}

func Service(request *http.Request, param map[string]string, reply *Reply) {
	reply.With("123" + param["name"])
	panic(errors.New("test error"))
}

func testService(request *http.Request, param map[string]string, reply *Reply) {
	reply.With(param["name"])
}
