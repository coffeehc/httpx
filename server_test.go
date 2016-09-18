// server_test
package web

import (
	"errors"
	"net/http"
	"testing"

	"github.com/coffeehc/logger"
	"github.com/coffeehc/utils"
)

func TestServer(t *testing.T) {
	logger.InitLogger()
	server := NewHttpServer(nil)
	server.Register("/a/{name}/123", GET, Service)
	server.Register("/a/123/{name}", GET, testService)
	server.AddLastFilter("/*", SimpleAccessLogFilter)
	server.Start()
	utils.WaitStop()
}

func Service(request *http.Request, param map[string]string, reply *Reply) {
	reply.With("123" + param["name"])
	panic(errors.New("test error"))
}

func testService(request *http.Request, param map[string]string, reply *Reply) {
	reply.With(param["name"])
}
