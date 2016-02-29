// server
package web

import (
	"net"
	"net/http"
	"strings"

	"github.com/coffeehc/logger"
)

type Server struct {
	httpServer *http.Server
	router     *router
	listener   net.Listener
	config     *ServerConfig
}

//创建一个Server,参数可以为空,默认使用0.0.0.0:8888
func NewServer(serverConfig *ServerConfig) *Server {
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
		ErrorLog:       logger.CreatLoggerAdapter(logger.LOGGER_LEVEL_ERROR, "", "", conf.HttpErrorLogout),
	}
	this.httpServer = server
	logger.Info("start HttpServer :%s", conf.getServerAddr())
	if conf.OpenTLS {
		go func() {
			err := server.ListenAndServeTLS(conf.certFile, conf.keyFile)
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

func (this *Server) serverHttpHandler(responseWriter http.ResponseWriter, request *http.Request) {
	request.URL.Path = strings.Replace(request.URL.Path, "//", "/", -1)
	httpReply := newHttpReply(responseWriter)
	this.router.filter.doFilter(request, httpReply)
	httpReply.finishReply()
}

func (server *Server) Regedit(path string, method HttpMethod, requestHandler RequestHandler) error {
	err := server.router.matcher.regeditAction(path, method, requestHandler)
	if err != nil {
		logger.Error("注册 Handler 失败:%s", err)
	}
	return err
}

func (server *Server) AddFilter(uriPattern string, actionFilter Filter) {
	server.router.addFilter(newServletStyleUriPatternMatcher(uriPattern), actionFilter)
}

func (server *Server) AddFilterWithRegex(uriPattern string, actionFilter Filter) {
	server.router.addFilter(newRegexUriPatternMatcher(uriPattern), actionFilter)
}
