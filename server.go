package httpx

import (
	"net"
	"net/http"

	"fmt"

	"crypto/tls"

	"context"
	"time"

	"git.xiagaogao.com/coffee/boot/errors"
	"git.xiagaogao.com/coffee/boot/logs"
	"go.uber.org/zap"
)

//Server the http server interface
type Server interface {
	Start() <-chan errors.Error
	Stop()
	GetServerAddress() string
	RegisterHandlerFunc(path string, method RequestMethod, handlerFunc http.HandlerFunc) errors.Error
	RegisterHandler(path string, method RequestMethod, handler http.Handler) errors.Error
	Register(path string, method RequestMethod, requestHandler RequestHandler) errors.Error

	AddFirstFilter(uriPattern string, actionFilter Filter)
	AddLastFilter(uriPattern string, actionFilter Filter)
	AddFilterWithRegex(uriPattern string, actionFilter Filter)

	AddRequestErrorHandler(code int, handler RequestErrorHandler) errors.Error
}

type serverImpl struct {
	httpServer *http.Server
	router     *router
	listener   net.Listener
	config     *Config
	errorService errors.Service
	logger *zap.Logger

}

//NewServer 创建一个Http Server,参数可以为空,默认使用0.0.0.0:8888
func NewServer(serverConfig *Config,errorService errors.Service,logger *zap.Logger) Server {
	errorService = errorService.NewService("httpserver")
	if serverConfig == nil {
		serverConfig = new(Config)
	}
	return &serverImpl{router: newRouter(errorService,logger), config: serverConfig,errorService:errorService,logger:logger}
}

func (s *serverImpl) Stop() {
	if s.httpServer != nil {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		s.httpServer.Shutdown(ctx)
	}
}

func (s *serverImpl) Start() <-chan errors.Error {
	s.logger.Debug("serverConfig", logs.F_ExtendData(s.config))
	s.router.matcher.sort()
	conf := s.config
	server := &http.Server{
		Addr:           conf.getServerAddr(),
		Handler:        http.HandlerFunc(s.serverHandler),
		ReadTimeout:    conf.getReadTimeout(),
		MaxHeaderBytes: conf.MaxHeaderBytes,
		TLSConfig:      conf.TLSConfig,
		TLSNextProto:   conf.TLSNextProto,
		ConnState:      conf.ConnState,
		IdleTimeout:    conf.getIdleTimeout(),
	}
	s.httpServer = server
	s.logger.Info("start HttpServer", logs.F_ExtendData(conf.getServerAddr()))
	errorSign := make(chan errors.Error, 1)
	listener, err := net.Listen("tcp", conf.getServerAddr())
	//TODO listen Option
	if err != nil {
		s.logger.Error("绑定监听地址失败", logs.F_ExtendData(conf.getServerAddr()))
		errorSign <- s.errorService.WappedSystemError(err)
		return errorSign
	}
	s.listener = &tcpKeepAliveListener{Listener: listener, keepAliveDuration: conf.getKeepAliveDuration()}
	if conf.TLSConfig != nil {
		server.TLSConfig = conf.TLSConfig
		s.listener = tls.NewListener(s.listener, server.TLSConfig)
	}
	go func() {
		err := server.Serve(s.listener)
		errorSign <- s.errorService.SystemError("启动 HttpServer 失败", logs.F_Error(err))
	}()
	return errorSign
}

func (s *serverImpl) GetServerAddress() string {
	return s.config.getServerAddr()
}

func (s *serverImpl) serverHandler(responseWriter http.ResponseWriter, request *http.Request) {
	reply := newHTTPReply(request, responseWriter, s.config,s.errorService,s.logger)
	defer func() {
		if err := recover(); err != nil {
			var httpErr *HTTPError
			var ok bool
			if httpErr, ok = err.(*HTTPError); !ok {
				httpErr = NewHTTPErr(500, fmt.Sprintf("%s", err))
			}
			s.logger.Error("http err", logs.F_ExtendData(err))
			reply.SetStatusCode(httpErr.Code)
			if handler, ok := s.router.errorHandlers[httpErr.Code]; ok {
				handler(httpErr, reply)
				return
			}
			reply.With(httpErr.Message).As(DefaultRenderJSON)
		}
		reply.finishReply()
	}()
	s.router.filter.doFilter(reply)

}

func (s *serverImpl) RegisterHandlerFunc(path string, method RequestMethod, handlerFunc http.HandlerFunc) errors.Error {
	return s.RegisterHandler(path, method, handlerFunc)
}

//适配 Http原生的 Handler 接口
func (s *serverImpl) RegisterHandler(path string, method RequestMethod, handler http.Handler) errors.Error {
	requestHandler := func(reply Reply) {
		reply.AdapterHTTPHandler(true)
		handler.ServeHTTP(reply.GetResponseWriter(), reply.GetRequest())
	}
	return s.Register(path, method, requestHandler)
}

func (s *serverImpl) Register(path string, method RequestMethod, requestHandler RequestHandler) errors.Error {
	return s.router.matcher.regeditAction(path, method, requestHandler)
}

func (s *serverImpl) AddFirstFilter(uriPattern string, actionFilter Filter) {
	s.router.addFirstFilter(newServletStyleURIPatternMatcher(uriPattern,s.logger), actionFilter)
}

func (s *serverImpl) AddLastFilter(uriPattern string, actionFilter Filter) {
	s.router.addLastFilter(newServletStyleURIPatternMatcher(uriPattern,s.logger), actionFilter)
}

func (s *serverImpl) AddFilterWithRegex(uriPattern string, actionFilter Filter) {
	s.router.addLastFilter(newRegexURIPatternMatcher(uriPattern,s.logger), actionFilter)
}

func (s *serverImpl) AddRequestErrorHandler(code int, handler RequestErrorHandler) errors.Error {
	if _, ok := s.router.errorHandlers[code]; ok {
		return s.errorService.SystemError(fmt.Sprintf("已经注册了[%d]异常响应码的处理方法,注册失败", code))
	}
	s.router.errorHandlers[code] = handler
	return nil
}
