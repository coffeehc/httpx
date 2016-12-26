package httpx

import (
	"net"
	"net/http"

	"fmt"

	"crypto/tls"
	"errors"
	"github.com/coffeehc/logger"
	"os"
)

//Server the http server interface
type Server interface {
	Start() <-chan error
	Stop()
	GetServerAddress() string
	RegisterHandlerFunc(path string, method RequestMethod, handlerFunc http.HandlerFunc) error
	RegisterHandler(path string, method RequestMethod, handler http.Handler) error
	Register(path string, method RequestMethod, requestHandler RequestHandler) error

	AddFirstFilter(uriPattern string, actionFilter Filter)
	AddLastFilter(uriPattern string, actionFilter Filter)
	AddFilterWithRegex(uriPattern string, actionFilter Filter)

	AddRequestErrorHandler(code int, handler RequestErrorHandler) error
}

type _Server struct {
	httpServer *http.Server
	router     *router
	listener   net.Listener
	config     *Config
}

//NewServer 创建一个Http Server,参数可以为空,默认使用0.0.0.0:8888
func NewServer(serverConfig *Config) Server {
	if serverConfig == nil {
		serverConfig = new(Config)
	}
	return &_Server{router: newRouter(), config: serverConfig}
}

func (s *_Server) Stop() {
	if s.listener != nil {
		logger.Info("http Listener Close")
		s.listener.Close()
	}
}

func (s *_Server) Start() <-chan error {
	logger.Debug("serverConfig is %#v", s.config)
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
		ErrorLog:       logger.CreatLoggerAdapter(logger.LoggerLevelError, "", "", os.Stdout),
	}
	if conf.HTTPErrorLogout != nil {
		server.ErrorLog = logger.CreatLoggerAdapter(logger.LoggerLevelError, "", "", conf.HTTPErrorLogout)
	}
	s.httpServer = server
	logger.Info("start HttpServer :%s", conf.getServerAddr())
	errorSign := make(chan error, 1)
	listener, err := net.Listen("tcp", conf.getServerAddr())
	//TODO listen Option
	if err != nil {
		logger.Error("绑定监听地址[%s]失败", conf.getServerAddr())
		errorSign <- err
		return errorSign
	}
	s.listener = &tcpKeepAliveListener{Listener: listener, keepAliveDuration: conf.getKeepAliveDuration()}
	if conf.TLSConfig != nil {
		server.TLSConfig = conf.TLSConfig
		s.listener = tls.NewListener(s.listener, server.TLSConfig)
	}
	go func() {
		err := server.Serve(s.listener)
		errorSign <- errors.New(logger.Error("启动 HttpServer 失败:%s", err))
	}()
	return errorSign
}

func (s *_Server) GetServerAddress() string {
	return s.config.getServerAddr()
}

func (s *_Server) serverHandler(responseWriter http.ResponseWriter, request *http.Request) {
	reply := newHTTPReply(request, responseWriter, s.config)
	defer func() {
		if err := recover(); err != nil {
			var httpErr *HTTPError
			var ok bool
			if httpErr, ok = err.(*HTTPError); !ok {
				httpErr = NewHTTPErr(500, fmt.Sprintf("%s", err))
			}
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

func (s *_Server) RegisterHandlerFunc(path string, method RequestMethod, handlerFunc http.HandlerFunc) error {
	return s.RegisterHandler(path, method, handlerFunc)
}

//适配 Http原生的 Handler 接口
func (s *_Server) RegisterHandler(path string, method RequestMethod, handler http.Handler) error {
	requestHandler := func(reply Reply) {
		reply.AdapterHTTPHandler(true)
		handler.ServeHTTP(reply.GetResponseWriter(), reply.GetRequest())
	}
	return s.Register(path, method, requestHandler)
}

func (s *_Server) Register(path string, method RequestMethod, requestHandler RequestHandler) error {
	err := s.router.matcher.regeditAction(path, method, requestHandler)
	if err != nil {
		logger.Error("注册 Handler 失败:%s", err)
	}
	return err
}

func (s *_Server) AddFirstFilter(uriPattern string, actionFilter Filter) {
	s.router.addFirstFilter(newServletStyleURIPatternMatcher(uriPattern), actionFilter)
}

func (s *_Server) AddLastFilter(uriPattern string, actionFilter Filter) {
	s.router.addLastFilter(newServletStyleURIPatternMatcher(uriPattern), actionFilter)
}

func (s *_Server) AddFilterWithRegex(uriPattern string, actionFilter Filter) {
	s.router.addLastFilter(newRegexURIPatternMatcher(uriPattern), actionFilter)
}

func (s *_Server) AddRequestErrorHandler(code int, handler RequestErrorHandler) error {
	if _, ok := s.router.errorHandlers[code]; ok {
		return errors.New(logger.Error("已经注册了[%d]异常响应码的处理方法,注册失败", code))
	}
	s.router.errorHandlers[code] = handler
	return nil
}
