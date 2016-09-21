// main
package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"os"

	"github.com/coffeehc/logger"
	"github.com/coffeehc/utils"
	"github.com/coffeehc/web"
	"github.com/coffeehc/web/pprof"
)

func main() {
	logger.InitLogger()
	config := &web.HttpServerConfig{
		EnabledTLS:      false,
		CertFile:        "server.crt",
		KeyFile:         "server.key",
		HttpErrorLogout: os.Stderr,
	}
	server := web.NewHttpServer(config)
	pprof.RegeditPprof(server)
	server.Register("/test", web.GET, TestService)
	server.Register("/reqinfo", web.GET, reqInfoHandler)
	server.Register("/a/{name}/123", web.GET, Service)
	server.Register("/a", web.GET, getStruct)
	server.AddFirstFilter("/*", web.SimpleAccessLogFilter)
	errSign := server.Start()
	go func() {
		err := <-errSign
		panic(err)
	}()
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
	responseType := reply.GetRequest().FormValue("type")
	if responseType == "xml" {
		reply.With(test).As(web.Default_Render_Xml)
		return
	}
	reply.With(test).As(web.Default_Render_Json)
}

func reqInfoHandler(reply web.Reply) {
	stream := reply.GetResponseWriter()
	request := reply.GetRequest()
	fmt.Fprintf(stream, "Method: %s\n", request.Method)
	fmt.Fprintf(stream, "Protocol: %s\n", request.Proto)
	fmt.Fprintf(stream, "Host: %s\n", request.Host)
	fmt.Fprintf(stream, "RemoteAddr: %s\n", request.RemoteAddr)
	fmt.Fprintf(stream, "RequestURI: %q\n", request.RequestURI)
	fmt.Fprintf(stream, "URL: %#v\n", request.URL)
	fmt.Fprintf(stream, "Body.ContentLength: %d (-1 means unknown)\n", request.ContentLength)
	fmt.Fprintf(stream, "Close: %v (relevant for HTTP/1 only)\n", request.Close)
	fmt.Fprintf(stream, "TLS: %#v\n", request.TLS)
	fmt.Fprintf(stream, "\nHeaders:\n")
	reply.With(stream).As(web.Default_Render_Text)
}

func Service(reply web.Reply) {
	pathFragment := reply.GetPathFragment()
	name, err := pathFragment.GetAsString("name")
	if err != nil {
		panic("123")
	}
	reply.With("123" + name)
	panic("123")
}

func TestService(reply web.Reply) {
	stream := reply.GetResponseWriter()
	clientGone := stream.(http.CloseNotifier).CloseNotify()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	fmt.Fprintf(stream, "# ~1KB of junk to force browsers to start rendering immediately: \n")
	io.WriteString(stream, strings.Repeat("# xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\n", 13))
	for {
		fmt.Fprintf(stream, "%v\n", time.Now())
		stream.(http.Flusher).Flush()
		select {
		case <-ticker.C:
		case <-clientGone:
			logger.Info("Client %v disconnected from the clock", reply.GetRequest().RemoteAddr)
			return
		}
	}
}
