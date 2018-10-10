package httpx

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Service interface {
	Start(onShutdown func()) <-chan error
	Shutdown() error
	GetGinEngine() *gin.Engine
	NewRouterGroup(prefix string) *gin.RouterGroup
	GetServerAddress() string
}

func NewService(name string, config *Config, logger *zap.Logger) Service {
	logger.Debug(fmt.Sprintf("[%s]HTTP服务器配置", name), zap.Any("http_config", config))
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	if config == nil {
		config = &Config{}
	}
	gin.EnableJsonDecoderUseNumber()
	engine := gin.New()
	engine.RedirectFixedPath = true
	l, err := Listen(config.getServerAddr())
	if err != nil {
		logger.Error(fmt.Sprintf("[%s]创建HTTP服务失败", name), zap.Error(err))
		return nil
	}
	config.ServerAddr = l.Addr().String()
	l.Close()
	server := &http.Server{
		TLSConfig:    config.TLSConfig,
		TLSNextProto: config.TLSNextProto,
		ReadTimeout:  config.getReadTimeout(),
		WriteTimeout: config.getWriteTimeout(),
		Addr:         config.getServerAddr(),
		Handler:      engine,
	}
	impl := &serviceImpl{
		name:   name,
		config: config,
		logger: logger,
		engine: engine,
		server: server,
	}
	return impl
}

type serviceImpl struct {
	name   string
	config *Config
	logger *zap.Logger
	engine *gin.Engine
	server *http.Server
}

func (impl *serviceImpl) NewRouterGroup(prefix string) *gin.RouterGroup {
	return impl.engine.Group(prefix)
}

func (impl *serviceImpl) GetServerAddress() string {
	return impl.config.getServerAddr()
}

func (impl *serviceImpl) Shutdown() error {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*30)
	err := impl.server.Shutdown(ctx)
	if err != nil {
		impl.logger.Error("关闭HttpServer失败", zap.Error(err))
	}
	return err
}

func (impl *serviceImpl) GetGinEngine() *gin.Engine {
	return impl.engine
}

func (impl *serviceImpl) Start(onShutdown func()) <-chan error {
	errorSign := make(chan error, 1)
	impl.server.RegisterOnShutdown(func() {
		impl.logger.Debug(fmt.Sprintf("[%s]HTTP服务器关闭", impl.name))
		if onShutdown != nil {
			onShutdown()
		}
	})
	go func() {
		err := impl.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			impl.logger.Error(fmt.Sprintf("[%s]HTTP服务异常关闭", impl.name), zap.Error(err))
		}
		impl.logger.Debug(fmt.Sprintf("[%s]HTTP服务关闭", impl.name))
		errorSign <- err
	}()
	impl.logger.Debug(fmt.Sprintf("[%s]HTTP服务启动", impl.name), zap.String("address", impl.config.getServerAddr()))
	return errorSign
}
