/**
 * Created by coffee on 15/11/15.
 */
package utils
import (
	"os"
	"syscall"
	"time"
	"os/signal"
	"github.com/coffeehc/logger"
)

type serviceWarp struct {
	runFunc  func() error
	stopFunc func() error
}

func (this *serviceWarp)Run() error {
	return this.runFunc()
}

func (this *serviceWarp)Stop() error {
	return this.stopFunc()
}
func NewService(runFunc func() error, stopFunc func() error) Service {
	return &serviceWarp{runFunc, stopFunc}
}

func StartService(service Service) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("service crash,cause is %s", err)
		}
		if service.Stop != nil {
			stopErr := service.Stop()
			if stopErr != nil {
				logger.Error("关闭服务失败,%s", stopErr)
			}
		}
		time.Sleep(time.Second)
	}()
	if service.Run == nil {
		panic("没有指定Run方法")
	}
	err := service.Run()
	if err != nil {
		panic(logger.Error("服务运行错误:%s", err))
	}
	logger.Info("服务已正常启动")
	WaitStop()
}

type Service interface {
	Run() error
	Stop() error
}

/*
	wait,一般是可执行函数的最后用于阻止程序退出
*/
func WaitStop() {
	var sigChan = make(chan os.Signal, 3)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	sig := <-sigChan
	logger.Debug("接收到指令:%s,立即关闭程序", sig)
}

