package main

import (
	"flag"
	"net"
	"os"
	"time"

	"github.com/coffeehc/logger"
	"github.com/coffeehc/utils"
	"github.com/coffeehc/web"
	"github.com/coffeehc/web/static"
)

var (
	port = flag.String("port", "8888", "static http server port,defaule 8888")
	addr = flag.String("addr", "0.0.0.0", "static http server addr,default 0.0.0.0")
	path = flag.String("path", "", "static path ,defaule ./")
)

func main() {
	flag.Parse()
	logger.SetDefaultLevel("/", logger.LOGGER_LEVEL_INFO)
	defer logger.WaitToClose()
	if *path == "" {
		_path, _ := os.Getwd()
		flag.Set("path", _path)
	}
	logger.Info("static dir is %s", *path)
	config := &web.ServerConfig{
		ServerAddr:     net.JoinHostPort(*addr, *port),
		ReadTimeout:    time.Minute * 5,
		WriteTimeout:   time.Minute * 5,
		MaxHeaderBytes: 100000,
	}
	server := web.NewServer(config)
	static.RegisterStaticFilter(server, "/*", *path)
	server.AddFirstFilter("/*", web.SimpleAccessLogFilter)
	server.Start()
	utils.WaitStop()
}
