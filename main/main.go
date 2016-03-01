// main
package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/coffeehc/logger"
	"github.com/coffeehc/utils"
	"github.com/coffeehc/web"
	"github.com/coffeehc/web/pprof"
	"os"
)

func main() {
	logger.InitLogger()
	config := &web.ServerConfig{
		OpenTLS:         true,
		CertFile:        "server.crt",
		KeyFile:         "server.key",
		HttpErrorLogout: os.Stderr,
	}
	server := web.NewServer(config)
	pprof.RegeditPprof(server)
	server.Regedit("/test", web.GET, TestService)
	server.Regedit("/reqinfo", web.GET, reqInfoHandler)
	server.Regedit("/a/{name}/123", web.GET, Service)
	server.AddFilter("/*", web.SimpleAccessLogFilter)
	server.Start()
	utils.WaitStop()
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
	panic(web.HTTPERR_400("testtest 400"))
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
