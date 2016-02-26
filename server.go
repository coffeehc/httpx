// server
package web

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/coffeehc/logger"
	"github.com/golang/net/http2"
)

type HttpMethod string

const (
	GET     = HttpMethod("GET")
	POST    = HttpMethod("POST")
	PUT     = HttpMethod("PUT")
	DELETE  = HttpMethod("DELETE")
	HEAD    = HttpMethod("HEAD")
	TRACE   = HttpMethod("TRACE")
	CONNECT = HttpMethod("CONNECT")
	PATCH   = HttpMethod("PATCH")
)

type ServerConfig struct {
	Addr           string
	Port           int
	ReadTimeout    time.Duration // 读的最大Timeout时间
	WriteTimeout   time.Duration // 写的最大Timeout时间
	MaxHeaderBytes int           // 请求头的最大长度
	TLSConfig      *tls.Config   // 配置TLS
	serverAddr     string
	OpenHttp2      bool //是否开启http2
}

type Server struct {
	router   *routingDispatcher
	listener net.Listener
	config   *ServerConfig
}

//创建一个Server,参数可以为空,默认使用0.0.0.0:8888
func NewServer(serverConfig *ServerConfig) *Server {
	if serverConfig == nil {
		serverConfig = &ServerConfig{Addr: "0.0.0.0"}
	}
	if serverConfig.Port == 0 {
		serverConfig.Port = 8888
	}
	if serverConfig.OpenHttp2 && serverConfig.TLSConfig == nil {
		logger.Error("open http2 need TLS support")
		return nil
	}
	addr := net.JoinHostPort(serverConfig.Addr, strconv.Itoa(serverConfig.Port))
	serverAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		logger.Error("can't parse the server addr [%s],cause:%s", addr, err)
		return nil
	}
	serverConfig.serverAddr = serverAddr.String()
	return &Server{router: newRoutingDispatcher(), config: serverConfig}
}

func (this *Server) Start() error {
	logger.Debug("serverConfig is %#v", this.config)
	this.router.matcher.sort()
	conf := this.config
	server := &http.Server{Handler: http.HandlerFunc(this.serverHttpHandler), MaxHeaderBytes: conf.MaxHeaderBytes, TLSConfig: conf.TLSConfig}
	if conf.ReadTimeout > 0 {
		server.ReadTimeout = conf.ReadTimeout
	}
	if conf.WriteTimeout > 0 {
		server.WriteTimeout = conf.WriteTimeout
	}
	http2.ConfigureServer(server, &http2.Server{})
	var err error
	this.listener, err = net.Listen("tcp", conf.serverAddr)
	if err != nil {
		return errors.New(logger.Error("listen [%s] fail:%s", conf.serverAddr, err))
	}
	logger.Info("start HttpServer :%s", conf.serverAddr)
	keepAliveListrener := tcpKeepAliveListener{this.listener.(*net.TCPListener)}
	if conf.TLSConfig != nil {
		conf.TLSConfig.NextProtos = append(conf.TLSConfig.NextProtos, "http/1.1")
		go server.Serve(tls.NewListener(keepAliveListrener, conf.TLSConfig))
	} else {
		go server.Serve(keepAliveListrener)
	}
	return nil
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(time.Minute)
	return tc, nil
}

func (this *Server) serverHttpHandler(responseWriter http.ResponseWriter, request *http.Request) {
	request.URL.Path = strings.Replace(request.URL.Path, "//", "/", -1)
	reply := newReply(responseWriter)
	this.router.filter.doFilter(request, reply)
	//TODO 处理异常的StatusCode
	if !reply.openStream {
		responseWriter.Header().Set("Connection", "close")
		reply.write()
	}
	request.Body.Close()
}

func (this *Server) Stop() {
	if this.listener != nil {
		logger.Debug("Close Http Server")
		this.listener.Close()
	}
}

func (server *Server) Regedit(path string, method HttpMethod, requestHandler RequestHandler) error {
	err := server.router.matcher.regeditAction(path, method, requestHandler)
	if err != nil {
		logger.Error("注册 Handler 失败:%s", err)
	}
	return err
}

func (server *Server) RegeditWebSocket(path string, service WebScoketHandler) error {
	adapter := &webScoketAdapter{service}
	err := server.router.matcher.regeditAction(path, GET, adapter.webScoketHandlerAdapter)
	if err != nil {
		logger.Error("注册 Handler 失败:%s", err)
	}
	return err
}

func (server *Server) AddFilter(uriPattern string, actionFilter ActionFilter) {
	server.router.addFilter(newServletStyleUriPatternMatcher(uriPattern), actionFilter)
}

func (server *Server) AddFilterWithRegex(uriPattern string, actionFilter ActionFilter) {
	server.router.addFilter(newRegexUriPatternMatcher(uriPattern), actionFilter)
}
