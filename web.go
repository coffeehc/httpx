// web project web.go
package web

import (
	"fmt"
	"inject"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/coffeehc/logger"
)

type WebConfig struct {
	inject.Injector
	Host               string
	Port               int
	Context            string //上下文
	TemplateDir        string
	staticResourcePath http.FileSystem
	Welcome            string
	filterDefinitions  []*filterDefinition
	IsProjectModel     bool //是否产品模式
	Fastcgi            bool
}

func (this *WebConfig) SetStaticResourcePath(path string) {
	this.staticResourcePath = http.Dir(path)
}

func (this *WebConfig) AddFilter(filter Filter) {
	if this.filterDefinitions == nil {
		this.filterDefinitions = make([]*filterDefinition, 0)
	}
	if filter == nil {
		panic("不能指定一个空的Filter")
	}
	fd := new(filterDefinition)
	var pattern = filter.GetPattern()
	r, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("Filter指定的UriPattern不能解析:%s", pattern))
	}
	fd.pattern = r
	fd.filter = filter
	this.filterDefinitions = append(this.filterDefinitions, fd)
}

type HttpServer struct {
	serverListener net.Listener
}

var dispatcher *routingDispatcher = newRoutingDispatcher()

func (this *WebConfig) initConfig() {
	if this.Injector == nil {
		this.Injector = inject.New()
	}
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	this.Binding(path, nil, Bind_Key_AppPath)
	this.Binding(this.staticResourcePath, nil, Bind_Key_StaticResource)
	this.Binding(this.Welcome, nil, Bind_Key_Welcome)
	this.Binding(this.IsProjectModel, nil, Bind_Key_IsProjectModel)
	if this.TemplateDir != "" {
		this.Binding(this.TemplateDir, nil, Bind_Key_TemplateDir)
		templateRender := initTemplate(this.TemplateDir)
		this.Binding(templateRender, "", Bind_Key_Template)
	}
}

type globalHandler struct {
	http.Handler
	conf *WebConfig
}

func (this *globalHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	filter := newFilterChainInvocation(this.conf)
	filter.Binding(NewResponseWriter(w), (*http.ResponseWriter)(nil), "")
	filter.Binding(req, (*http.Request)(nil), "")
	reply := NewReply(w, filter.Injector)
	filter.Injector.Binding(reply, (*Reply)(nil), "")
	filter.DoFilter(req, reply)
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
	logger.Info("开始启动Http服务器")
	if conf.Port <= 0 || conf.Port >= 65536 {
		logger.Error("端口号不符合要求:%d", conf.Port)
		Stop()
	}
	conf.initConfig()
	for _, filter := range conf.filterDefinitions {
		filter.filter.Init(conf)
	}
	if conf.Host == "" {
		conf.Host = "0.0.0.0"
	}
	addr := net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port))
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("创建监听端口出现一个错误:%v", err)
	}
	httpServer.serverListener = l
	var handler = newGlobalHandler(conf)
	if conf.Fastcgi {
		fcgi.Serve(l, handler)
	} else {
		server := &http.Server{Addr: addr, Handler: handler}
		err = server.Serve(l)
	}

	if err != nil {
		return nil, fmt.Errorf("启动服务器出现一个错误:%v", err)
	}
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
