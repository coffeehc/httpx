// server
package web

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/coffeehc/logger"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	HEAD    = "HEAD"
	TRACE   = "TRACE"
	CONNECT = "CONNECT"
)

type defauleAction struct {
	path    string
	method  string
	service func(request *http.Request, pathFragments map[string]string, reply *Reply)
}

func (this *defauleAction) GetPath() string {
	return this.path
}
func (this *defauleAction) GetMethod() string {
	return this.method
}
func (this *defauleAction) Service(request *http.Request, pathFragments map[string]string, reply *Reply) {
	this.service(request, pathFragments, reply)
}

type ServerConfig struct {
	Addr           string
	Port           int
	ReadTimeout    time.Duration // 读的最大Timeout时间
	WriteTimeout   time.Duration // 写的最大Timeout时间
	MaxHeaderBytes int           // 请求头的最大长度
	TLSConfig      *tls.Config   // 配置TLS
	serverAddr     *net.TCPAddr
}

type Server struct {
	router   *routingDispatcher
	listener net.Listener
	config   *ServerConfig
}

//创建一个Server,参数可以为空,默认使用0.0.0.0:8888
func NewServer(serverConfig *ServerConfig) *Server {
	if serverConfig == nil {
		serverConfig = &ServerConfig{Addr: "0.0.0.0", Port: 8888}
	}
	addr := net.JoinHostPort(serverConfig.Addr, strconv.Itoa(serverConfig.Port))
	serverAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		logger.Error("设置的服务器地址[%s]无法解析:%s", addr, err)
		return nil
	}
	serverConfig.serverAddr = serverAddr
	return &Server{router: newRoutingDispatcher(), config: serverConfig}
}

func (this *Server) Start() error {
	conf := this.config
	server := &http.Server{Handler: http.HandlerFunc(this.serverHttpHandler), ReadTimeout: conf.ReadTimeout, WriteTimeout: conf.WriteTimeout, MaxHeaderBytes: conf.MaxHeaderBytes, TLSConfig: conf.TLSConfig}
	var err error
	this.listener, err = net.ListenTCP("tcp", conf.serverAddr)
	if err != nil {
		return errors.New(logger.Error("监听地址[%s]失败:%s", conf.serverAddr, err))
	}
	go server.Serve(this.listener)
	return nil
}

func (this *Server) serverHttpHandler(responseWriter http.ResponseWriter, request *http.Request) {
	reply := newReply()
	this.router.filters[0].filter(request, reply)
	//TODO 处理异常的StatusCode
	reply.write(responseWriter)
}

func (this *Server) Stop() {
	if this.listener != nil {
		logger.Debug("Close Http Server")
		this.listener.Close()
	}
}

func (server *Server) Regedit(path string, method string, service func(request *http.Request, pathFragments map[string]string, reply *Reply)) error {
	return server.router.matcher.regeditAction(&defauleAction{path, method, service})
}

func (server *Server) AddFilter(uriPattern string, actionFilter ActionFilter) {
	server.router.addFilter(newServletStyleUriPatternMatcher(uriPattern), actionFilter)
}

func (server *Server) AddFilterWithRegex(uriPattern string, actionFilter ActionFilter) {
	server.router.addFilter(newRegexUriPatternMatcher(uriPattern), actionFilter)
}
