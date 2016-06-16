// main
package web_test

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

func ExampleMain() {
	logger.InitLogger()
	config := &web.ServerConfig{
		OpenTLS:         true,
		CertFile:        "server.crt",
		KeyFile:         "server.key",
		HttpErrorLogout: os.Stderr,
	}
	server := web.NewServer(config)
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

func getStruct(r *http.Request, param map[string]string, reply web.Reply) {
	test := &Test{
		Name: "coffee",
		Age:  1,
		Sex:  0,
	}
	responseType := r.FormValue("type")
	if responseType == "xml" {
		reply.With(test).As(web.Transport_Xml)
		return
	}
	reply.With(test).As(web.Transport_Json)
}

func reqInfoHandler(r *http.Request, param map[string]string, reply web.Reply) {
	stream := reply.GetResponseWriter()
	fmt.Fprintf(stream, "Method: %s\n", r.Method)
	fmt.Fprintf(stream, "Protocol: %s\n", r.Proto)
	fmt.Fprintf(stream, "Host: %s\n", r.Host)
	fmt.Fprintf(stream, "RemoteAddr: %s\n", r.RemoteAddr)
	fmt.Fprintf(stream, "RequestURI: %q\n", r.RequestURI)
	fmt.Fprintf(stream, "URL: %#v\n", r.URL)
	fmt.Fprintf(stream, "Body.ContentLength: %d (-1 means unknown)\n", r.ContentLength)
	fmt.Fprintf(stream, "Close: %v (relevant for HTTP/1 only)\n", r.Close)
	fmt.Fprintf(stream, "TLS: %#v\n", r.TLS)
	fmt.Fprintf(stream, "\nHeaders:\n")
	r.Header.Write(stream)
}

func Service(request *http.Request, param map[string]string, reply web.Reply) {
	reply.With("123" + param["name"])
	panic("123")
}

func TestService(request *http.Request, param map[string]string, reply web.Reply) {
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
			logger.Info("Client %v disconnected from the clock", request.RemoteAddr)
			return
		}
	}
}
