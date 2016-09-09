// server
package web

import (
	"net"
	"net/http"

	"fmt"

	"errors"

	"github.com/coffeehc/logger"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/valyala/fasthttp"
)

type HttpServer interface {
	Start() error
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
	httpServer *fasthttp.Server
	router     *router
	listener   net.Listener
	config     *HttpServerConfig
}

//创建一个Server,参数可以为空,默认使用0.0.0.0:8888
func NewServer(serverConfig *HttpServerConfig) HttpServer {
	if serverConfig == nil {
		serverConfig = new(HttpServerConfig)
	}
	return &_Server{router: newRouter(), config: serverConfig}
}

func (this *_Server) Start() error {
	logger.Debug("serverConfig is %#v", this.config)
	this.router.matcher.sort()
	conf := this.config
	option := this.config.GetServerOption()
	server := &fasthttp.Server{
		Handler:                       this.httpHandler,
		Name:                          conf.GetName(),
		Concurrency:                   option.GetConcurrency(),
		DisableKeepalive:              option.GetDisableKeepalive(),
		ReadBufferSize:                option.GetReadBufferSize(),
		WriteBufferSize:               option.GetWriteBufferSize(),
		ReadTimeout:                   option.GetReadTimeout(),
		WriteTimeout:                  option.GetWriteTimeout(),
		MaxConnsPerIP:                 option.GetMaxConnsPerIP(),
		MaxRequestsPerConn:            option.GetMaxRequestsPerConn(),
		MaxRequestBodySize:            option.GetMaxRequestBodySize(),
		ReduceMemoryUsage:             option.GetReduceMemoryUsage(),
		GetOnly:                       false,
		LogAllErrors:                  true,
		DisableHeaderNamesNormalizing: false,
		Logger: httpLogger{},
	}
	this.httpServer = server
	logger.Info("start HttpServer :%s", conf.GetServerAddr())
	if conf.OpenTLS {
		go func() {
			err := server.ListenAndServeTLS(conf.GetServerAddr(), conf.CertFile, conf.KeyFile)
			logger.Error("启动 HttpServer 失败:%s", err)
		}()
	} else {
		go func() {
			err := server.ListenAndServe(conf.GetServerAddr())
			logger.Error("启动 HttpServer 失败:%s", err)
		}()
	}
	return nil
}

func (this *_Server) GetServerAddress() string {
	return this.config.GetServerAddr()
}

func (this *_Server) httpHandler(ctx *fasthttp.RequestCtx) {
	//TODO 考虑 context 使用 timeoutContext
	reply := newReply(ctx, context.TODO(), this.config.GetDefaultRender())
	ctx.Response.SetStatusCode(200)
	defer func() {
		if err := recover(); err != nil {
			var httpErr *HttpError
			var ok bool
			if httpErr, ok = err.(*HttpError); !ok {
				httpErr = HTTPERR_500(fmt.Sprintf("%#s", err))
			}
			defer reply.SetStatusCode(httpErr.Code)
			if handler, ok := this.router.errorHandlers[httpErr.Code]; ok {
				handler(httpErr, reply)
				return
			}
			reply.With(httpErr.Message).As(Render_Json)
		}
		err := reply.FinishReply()
		if err != nil {
			ctx.Response.SetStatusCode(500)
			ctx.Response.SetBodyString(err.Error())
		}
	}()
	this.router.filter.doFilter(reply)
}

func (server *_Server) RegisterHttpHandlerFunc(path string, method HttpMethod, handlerFunc http.HandlerFunc) error {
	return server.RegisterHttpHandler(path, method, handlerFunc)
}

//适配 Http原生的 Handler 接口
func (server *_Server) RegisterHttpHandler(path string, method HttpMethod, handler http.Handler) error {
	//requestHandler := func(request *http.Request, pathFragments map[string]string, reply Reply) {
	//	reply.AdapterHttpHandler(true)
	//	handler.ServeHTTP(reply.GetResponseWriter(), request)
	//}
	//return server.Register(path, method, requestHandler)
	//TODO 未完成
	return nil
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
		return errors.New(logger.Error("已经注册了[%d]异常响应码的处理方法,注册失败"))
	}
	server.router.errorHandlers[code] = handler
	return nil
}
