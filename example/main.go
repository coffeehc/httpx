// main
package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/coffeehc/commons"
	"github.com/coffeehc/httpx"
	"github.com/coffeehc/httpx/pprof"
	"github.com/coffeehc/logger"
)

func main() {
	var mu sync.Mutex
	var items = make(map[int]struct{})

	runtime.SetMutexProfileFraction(5)
	for i := 0; i < 1000*1000; i++ {
		go func(i int) {
			mu.Lock()
			defer mu.Unlock()
			items[i] = struct{}{}
		}(i)
	}

	http.ListenAndServe(":8888", nil)
}

func main1() {
	logger.InitLogger()
	defer logger.WaitToClose()
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("load cert err %s", err)
	}
	config := &httpx.Config{
		TLSConfig:       &tls.Config{Certificates: []tls.Certificate{cert}},
		HTTPErrorLogout: os.Stderr,
	}
	server := httpx.NewServer(config)
	pprof.RegeditPprof(server)
	server.Register("/test", httpx.GET, TestService)
	server.Register("/reqinfo", httpx.GET, reqInfoHandler)
	server.Register("/a/{name}/123", httpx.GET, Service)
	server.Register("/a", httpx.GET, getStruct)
	server.AddFirstFilter("/*", httpx.AccessLogFilter)
	errSign := server.Start()
	go func() {
		err := <-errSign
		panic(err)
	}()
	commons.WaitStop()
}

//Test test
type Test struct {
	Name string
	Age  int
	Sex  int
}

func getStruct(reply httpx.Reply) {
	test := &Test{
		Name: "coffee",
		Age:  1,
		Sex:  0,
	}
	responseType := reply.GetRequest().FormValue("type")
	if responseType == "xml" {
		reply.With(test).As(httpx.DefaultRenderXML)
		return
	}
	reply.With(test).As(httpx.DefaultRenderJSON)
}

func reqInfoHandler(reply httpx.Reply) {
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
	reply.With(stream).As(httpx.DefaultRenderText)
}

//Service Service
func Service(reply httpx.Reply) {
	pathFragment := reply.GetPathFragment()
	name, err := pathFragment.GetAsString("name")
	if err != nil {
		panic("123")
	}
	reply.With("123" + name)
	panic("123")
}

//TestService test Service
func TestService(reply httpx.Reply) {
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
