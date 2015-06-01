// main
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/coffeehc/logger"
	"github.com/coffeehc/web"
)

func main() {
	tlsConfig, err := web.NewTLSConfig("server.crt", "server.key")
	if err != nil {
		logger.Error("证书初始化失败,%s", err)
		time.Sleep(time.Second)
		return
	}
	serConfig := new(web.ServerConfig)
	serConfig.TLSConfig = tlsConfig
	serConfig.OpenHttp2 = true
	server := web.NewServer(serConfig)
	server.Regedit("/test", web.GET, TestService)
	server.Regedit("/reqinfo", web.GET, reqInfoHandler)
	server.Regedit("/a/{name}/123", web.GET, Service)
	server.AddFilter("/*", web.AccessLogFilter)
	server.Start()
	time.Sleep(time.Second * 120)
	server.Stop()
}

func reqInfoHandler(r *http.Request, param map[string]string, reply *web.Reply) {
	stream := reply.OpenStream()
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

func Service(request *http.Request, param map[string]string, reply *web.Reply) {
	reply.With("123" + param["name"])
	panic(errors.New("test error"))
}

func TestService(request *http.Request, param map[string]string, reply *web.Reply) {
	stream := reply.OpenStream()
	clientGone := stream.CloseNotify()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	fmt.Fprintf(stream, "# ~1KB of junk to force browsers to start rendering immediately: \n")
	io.WriteString(stream, strings.Repeat("# xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\n", 13))
	for {
		fmt.Fprintf(stream, "%v\n", time.Now())
		stream.Flush()
		select {
		case <-ticker.C:
		case <-clientGone:
			log.Printf("Client %v disconnected from the clock", request.RemoteAddr)
			return
		}
	}
}
