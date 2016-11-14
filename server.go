// server
package web

import (
	"net"
	"net/http"
	"strings"

	"fmt"

	"crypto/tls"
	"errors"
	"github.com/coffeehc/logger"
)

type HttpServer interface {
	Start() <-chan error
	Stop()
	GetServerAddress() string
	RegisterHttpHandlerFunc(path string, method HttpMethod, handlerFunc http.HandlerFunc) error
	RegisterHttpHandler(path string, method HttpMethod, handler http.Handler) error
	Register(path string, method HttpMethod, requestHandler RequestHandler) error

	AddFirstFilter(uriPattern string, actionFilter Filter)
	AddLastFilter(uriPattern string, actionFilter Filter)
	AddFilterWithRegex(uriPattern string, actionFilter Filter)

	AddRequestErrorHandler(code int, handler RequestErrorHandler) error
}

type _Server struct {
	httpServer *http.Server
	router     *router
	listener   net.Listener
	config     *HttpServerConfig
}

//创建一个Server,参数可以为空,默认使用0.0.0.0:8888
func NewHttpServer(serverConfig *HttpServerConfig) HttpServer {
	if serverConfig == nil {
		serverConfig = new(HttpServerConfig)
	}
	return &_Server{router: newRouter(), config: serverConfig}
}

func (this *_Server) Stop() {
	if this.listener != nil {
		this.listener.Close()
	}
}

func (this *_Server) Start() <-chan error {
	logger.Debug("serverConfig is %#v", this.config)
	this.router.matcher.sort()
	conf := this.config
	server := &http.Server{
		Addr:           conf.getServerAddr(),
		Handler:        http.HandlerFunc(this.serverHttpHandler),
		ReadTimeout:    conf.getReadTimeout(),
		MaxHeaderBytes: conf.MaxHeaderBytes,
		TLSConfig:      conf.TLSConfig,
		TLSNextProto:   conf.TLSNextProto,
		ConnState:      conf.ConnState,
	}
	if conf.HttpErrorLogout != nil {
		server.ErrorLog = logger.CreatLoggerAdapter(logger.LOGGER_LEVEL_ERROR, "", "", conf.HttpErrorLogout)
	}
	this.httpServer = server
	logger.Info("start HttpServer :%s", conf.getServerAddr())
	errorSign := make(chan error, 1)
	listen, err := net.Listen("tcp", conf.getServerAddr())
	listen = tcpKeepAliveListener{TCPListener: listen.(*net.TCPListener), keepAliveDuration: conf.getKeepAliveDuration()}
	//TODO listen Option
	if err != nil {
		logger.Error("绑定监听地址[%s]失败", conf.getServerAddr())
		errorSign <- err
		return errorSign
	}
	if conf.EnabledTLS {
		cer, err := tls.LoadX509KeyPair(conf.CertFile, conf.KeyFile)
		if err != nil {
			logger.Error("加载 TSL 证书失败,%s", err)
			errorSign <- err
			return errorSign
		}
		if conf.TLSConfig == nil {
			tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
			server.TLSConfig = tlsConfig
		} else {
			server.TLSConfig.Certificates = []tls.Certificate{cer}
		}
		listen = tls.NewListener(listen, server.TLSConfig)
	}
	go func() {
		err := server.Serve(listen)
		errorSign <- errors.New(logger.Error("启动 HttpServer 失败:%s", err))
	}()
	return errorSign
}

func (this *_Server) GetServerAddress() string {
	return this.config.getServerAddr()
}

func (this *_Server) serverHttpHandler(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	request.URL.Path = strings.Replace(request.URL.Path, "//", "/", -1)
	reply := newHttpReply(request, responseWriter, this.config)
	defer func() {
		if err := recover(); err != nil {
			var httpErr *HttpError
			var ok bool
			if httpErr, ok = err.(*HttpError); !ok {
				httpErr = HTTPERR_500(fmt.Sprintf("%#s", err))
			}
			reply.SetStatusCode(httpErr.Code)
			if handler, ok := this.router.errorHandlers[httpErr.Code]; ok {
				handler(httpErr, reply)
				return
			}
			reply.With(httpErr.Message).As(Default_Render_Json)
		}
		reply.finishReply()
	}()
	this.router.filter.doFilter(reply)

}

func (server *_Server) RegisterHttpHandlerFunc(path string, method HttpMethod, handlerFunc http.HandlerFunc) error {
	return server.RegisterHttpHandler(path, method, handlerFunc)
}

//适配 Http原生的 Handler 接口
func (server *_Server) RegisterHttpHandler(path string, method HttpMethod, handler http.Handler) error {
	requestHandler := func(reply Reply) {
		reply.AdapterHttpHandler(true)
		handler.ServeHTTP(reply.GetResponseWriter(), reply.GetRequest())
	}
	return server.Register(path, method, requestHandler)
}

func (server *_Server) Register(path string, method HttpMethod, requestHandler RequestHandler) error {
	err := server.router.matcher.regeditAction(path, method, requestHandler)
	if err != nil {
		logger.Error("注册 Handler 失败:%s", err)
	}
	return err
}

func (server *_Server) AddFirstFilter(uriPattern string, actionFilter Filter) {
	server.router.addFirstFilter(newServletStyleUriPatternMatcher(uriPattern), actionFilter)
}

func (server *_Server) AddLastFilter(uriPattern string, actionFilter Filter) {
	server.router.addLastFilter(newServletStyleUriPatternMatcher(uriPattern), actionFilter)
}

func (server *_Server) AddFilterWithRegex(uriPattern string, actionFilter Filter) {
	server.router.addLastFilter(newRegexUriPatternMatcher(uriPattern), actionFilter)
}

func (server *_Server) AddRequestErrorHandler(code int, handler RequestErrorHandler) error {
	if _, ok := server.router.errorHandlers[code]; ok {
		return errors.New(logger.Error("已经注册了[%d]异常响应码的处理方法,注册失败", code))
	}
	server.router.errorHandlers[code] = handler
	return nil
}
