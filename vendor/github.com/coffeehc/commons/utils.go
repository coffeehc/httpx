package commons

import (
	"fmt"
	"github.com/coffeehc/logger"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

var (
	_localIP = net.IPv4(127, 0, 0, 1)
)

//GetLocalIP 获取本地Ip
func GetLocalIP() net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.Error("无法获取网络接口信息,%s", err)
		return _localIP
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP
			}
		}
	}
	return _localIP
}

//GetAppPath 获取 App 路径
func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	return path
}

//GetAppDir 获取App执行文件目录
func GetAppDir() string {
	return filepath.Dir(GetAppPath())
}

//WaitStop 一般是可执行函数的最后用于阻止程序退出
func WaitStop() {
	var sigChan = make(chan os.Signal, 3)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	sig := <-sigChan
	fmt.Printf("接收到指令:%s,立即关闭程序", sig)
}
