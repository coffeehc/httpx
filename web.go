// web project web.go
package web

import "net/http"
import "logger"
import "time"
import "os"
import "net"
import "strconv"
import "fmt"

type WebConfig struct {
	StaticDir string
	Host      string
	Port      int
	isDev     bool
	context   string //上下文
}

var fileServer http.Handler

var dispatcher *routingDispatcher

func init() {
	logger.Debug("初始化WebServer")
	dispatcher = newRoutingDispatcher()
}

func (this *WebConfig) initConfig() {
	logger.Info("初始化Web配置,未做任何事")
	//if strings.HasPrefix(this.StaticDir,"/")
}

type globalHandler struct {
	http.Handler
	conf *WebConfig
}

func (this *globalHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger.Debugf("接收到一个请求,path:%s", req.RequestURI)
	logger.Debug("开始构建request和Response")
	filter := newFilterChainInvocation(this.conf)
	//TODO 此处是否需要处理一下RequseWap
	err := filter.doFilter(req, w)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "出现了一个错误:\n%s", err)
	}
}

func newGlobalHandler(conf *WebConfig) *globalHandler {
	//此处留作初始化一些信息用的
	handler := new(globalHandler)
	handler.conf = conf
	return handler
}

//启动服务器
func Strat(conf *WebConfig) {
	logger.Info("启动服务器")
	if conf.Port <= 0 || conf.Port >= 65536 {
		logger.Errorf("端口号不符合要求:%d", conf.Port)
		Stop()
	}
	conf.initConfig()
	initRoute(conf)
	fileServer = http.FileServer(http.Dir(conf.StaticDir))
	err := http.ListenAndServe(net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port)), newGlobalHandler(conf))
	if err != nil {
		logger.Errorf("启动服务器出现一个错误:%v", err)
		Stop()
	}
}

func Stop() {
	time.Sleep(time.Second * 3)
	os.Exit(80)
}

func initRoute(conf *WebConfig) {

}
