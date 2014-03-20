// web project web.go
package web

import "net/http"
import "github.com/coffeehc/logger"
import "time"
import "os"
import "net"
import "strconv"
import "fmt"

type WebConfig struct {
	StaticDir string
	Host      string
	Port      int
	Context   string //上下文
	Welcome   string
}
type HttpServer struct {
	serverListener net.Listener
}

var dispatcher *routingDispatcher
var fileServer fileHandler

func init() {
	logger.Debug("初始化WebServer")
	dispatcher = newRoutingDispatcher()
}

func (this *WebConfig) initConfig() {
	if this.StaticDir != "" {
		fileServer.root = http.Dir(this.StaticDir)
	}
}

type globalHandler struct {
	http.Handler
	conf *WebConfig
}

func (this *globalHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	filter := newFilterChainInvocation(this.conf)
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("出现了不可处理的错误:%s", err)
		}
	}()
	reply := NewReply(w)
	err := filter.DoFilter(req, reply)
	if err != nil {
		reply.Error(fmt.Sprintf("出现了不可处理的异常:\n%s\n", err), 500)
	}
	reply.writeResponse(w, req)
}

func newGlobalHandler(conf *WebConfig) *globalHandler {
	//此处留作初始化一些信息用的
	handler := new(globalHandler)
	handler.conf = conf
	return handler
}

//启动服务器
func Strat(conf *WebConfig) (*HttpServer, error) {
	httpServer := new(HttpServer)
	logger.Info("启动服务器")
	if conf.Port <= 0 || conf.Port >= 65536 {
		logger.Errorf("端口号不符合要求:%d", conf.Port)
		Stop()
	}
	conf.initConfig()
	httpServer.initRoute(conf)
	server := &http.Server{Addr: net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port)), Handler: newGlobalHandler(conf)}
	addr := server.Addr
	if addr == "" {
		addr = ":http"
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("创建监听端口出现一个错误:%v", err)
	}
	err = server.Serve(l)
	if err != nil {
		return nil, fmt.Errorf("启动服务器出现一个错误:%v", err)
	}
	httpServer.serverListener = l
	return httpServer, nil
}

func (this *HttpServer) Close() error {
	if this.serverListener != nil {
		err := this.serverListener.Close()
		if err != nil {
			return fmt.Errorf("关闭http监听出现错误:%s", err)
		}
	}
	return nil
}

func Stop() {
	time.Sleep(time.Second * 3)
	os.Exit(80)
}

func (this *HttpServer) initRoute(conf *WebConfig) {

}
