// main
package main

import (
	"fmt"
	"io"
	"strings"

	"bytes"
	"github.com/coffeehc/logger"
	"github.com/coffeehc/utils"
	"github.com/coffeehc/web"
	"github.com/coffeehc/web/pprof"
	"time"
)

func main() {
	logger.InitLogger()
	config := &web.HttpServerConfig{
		OpenTLS:  false,
		CertFile: "server.crt",
		KeyFile:  "server.key",
		ServerOption: &web.ServerOption{
			WriteBufferSize: 32,
			WriteTimeout:    3,
		},
	}
	server := web.NewHttpServer(config)
	pprof.RegeditPprof(server)
	server.Register("/test", web.GET, TestService)
	server.Register("/reqinfo", web.GET, reqInfoHandler)
	server.Register("/a/{name}/123", web.GET, Service)
	server.Register("/a", web.GET, getStruct)
	server.AddFirstFilter("/*", web.SimpleAccessLogFilter)
	server.Start()
	utils.WaitStop()
}

type Test struct {
	Name string
	Age  int
	Sex  int
}

func getStruct(reply web.Reply) {
	test := &Test{
		Name: "coffee",
		Age:  1,
		Sex:  0,
	}
	responseType := reply.GetPostParam("type")
	if responseType.AsString() == "xml" {
		reply.With(test).As(web.Render_Text)
		return
	}
	reply.With(test).As(web.Render_Json)
}

func reqInfoHandler(reply web.Reply) {
	buf := bytes.NewBuffer(nil)
	cxt := reply.GetRequestContext()
	fmt.Fprintf(buf, "Method: %s\n", reply.GetHttpMethod())
	fmt.Fprintf(buf, "Protocol http1.1: %t\n", cxt.Request.Header.IsHTTP11())
	fmt.Fprintf(buf, "Host: %s\n", cxt.Host())
	fmt.Fprintf(buf, "RemoteAddr: %s\n", cxt.RemoteAddr())
	fmt.Fprintf(buf, "RequestURI: %s\n", cxt.RequestURI())
	fmt.Fprintf(buf, "URL: %#v\n", cxt.Request.URI())
	fmt.Fprintf(buf, "Body.ContentLength: %d (-1 means unknown)\n", cxt.Request.Header.ContentLength())
	//fmt.Fprintf(buf, "Close: %v (relevant for HTTP/1 only)\n", r.Close)
	//fmt.Fprintf(buf, "TLS: %#v\n", r.TLS)
	fmt.Fprintf(buf, "\nHeaders:\n")
	reply.With(buf.Bytes())
}

func Service(reply web.Reply) {
	logger.Debug("StatusCode is [%d]", reply.GetStatusCode())
	pathFragment := reply.GetPathFragment()
	reply.With("123" + pathFragment["name"].AsString())
	panic("123")
}

func TestService(reply web.Reply) {
	//stream := reply.GetResponseWriter()
	pipeR, pipeW := io.Pipe()
	go func() {
		fmt.Fprintf(pipeW, "# ~1KB of junk to force browsers to start rendering immediately: \n")
		io.WriteString(pipeW, strings.Repeat("# xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\n", 13))
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		remoteAddr := reply.GetRemoteAddr()
		for t := range ticker.C {
			_, err := fmt.Fprintf(pipeW, "%v\n", t)
			if err != nil {
				logger.Info("Client %v disconnected from the clock,error is %s", remoteAddr, err)
				return
			}
		}
	}()
	reply.With(pipeR)
}
