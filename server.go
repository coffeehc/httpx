// server
package web

import (
	"net"
	"net/http"
	"strings"

	"fmt"
	"github.com/coffeehc/logger"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

type Server struct {
	httpServer *http.Server
	router     *router
	listener   net.Listener
	config     *ServerConfig
}

//创建一个Server,参数可以为空,默认使用0.0.0.0:8888
func NewServer(serverConfig *ServerConfig) *Server {
	if serverConfig == nil {
		serverConfig = new(ServerConfig)
	}
	return &Server{router: newRouter(), config: serverConfig}
}

func (this *Server) Start() error {
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
	if conf.OpenTLS {
		go func() {
			err := server.ListenAndServeTLS(conf.CertFile, conf.KeyFile)
			logger.Error("启动 HttpServer 失败:%s", err)
		}()
	} else {
		go func() {
			err := server.ListenAndServe()
			logger.Error("启动 HttpServer 失败:%s", err)
		}()
	}
	return nil
}

func (this *Server) GetServerAddress() string {
	return this.config.getServerAddr()
}

func (this *Server) serverHttpHandler(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	request.URL.Path = strings.Replace(request.URL.Path, "//", "/", -1)
	reply := newHttpReply(responseWriter, this.config)
	defer func() {
		if err := recover(); err != nil {
			var httpErr *HttpError
			var ok bool
			if httpErr, ok = err.(*HttpError); !ok {
				httpErr = HTTPERR_500(fmt.Sprintf("%#s", err))
			}
			defer reply.SetStatusCode(httpErr.Code)
			if handler, ok := this.router.errorHandlers[httpErr.Code]; ok {
				handler(request, httpErr, reply)
				return
			}
			reply.With(httpErr.Message).As(Transport_Json)
		}
		reply.finishReply(request, this.config.GetRender())
	}()
	this.router.filter.doFilter(request, reply)

}

func (server *Server) RegisterHttpHandlerFunc(path string, method HttpMethod, handlerFunc http.HandlerFunc) error {
	return server.RegisterHttpHandler(path, method, handlerFunc)
}

//适配 Http原生的 Handler 接口
func (server *Server) RegisterHttpHandler(path string, method HttpMethod, handler http.Handler) error {
	requestHandler := func(request *http.Request, pathFragments map[string]string, reply Reply) {
		reply.AdapterHttpHander(true)
		handler.ServeHTTP(reply.GetResponseWriter(), request)
	}
	return server.Register(path, method, requestHandler)
}

func (server *Server) Register(path string, method HttpMethod, requestHandler RequestHandler) error {
	err := server.router.matcher.regeditAction(path, method, requestHandler)
	if err != nil {
		logger.Error("注册 Handler 失败:%s", err)
	}
	return err
}

func (server *Server) AddFirstFilter(uriPattern string, actionFilter Filter) {
	server.router.addFirstFilter(newServletStyleUriPatternMatcher(uriPattern), actionFilter)
}

func (server *Server) AddLastFilter(uriPattern string, actionFilter Filter) {
	server.router.addLastFilter(newServletStyleUriPatternMatcher(uriPattern), actionFilter)
}

func (server *Server) AddFilterWithRegex(uriPattern string, actionFilter Filter) {
	server.router.addLastFilter(newRegexUriPatternMatcher(uriPattern), actionFilter)
}

func (server *Server) AddRequestErrorHandler(code int, handler RequestErrorHandler) error {
	if _, ok := server.router.errorHandlers[code]; ok {
		return errors.New(logger.Error("已经注册了[%d]异常响应码的处理方法,注册失败"))
	}
	server.router.errorHandlers[code] = handler
	return nil
}
