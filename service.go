package httpx

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Service interface {
	ReleaseMode()
	Start(onShutdown func()) <-chan error
	Shutdown() error
	GetGinEngine() *gin.Engine
	NewRouterGroup(prefix string)*gin.RouterGroup
}

func NewService(config *Config, logger *zap.Logger) (Service, error) {
	engine := gin.New()
	server := &http.Server{
		TLSConfig:    config.TLSConfig,
		TLSNextProto: config.TLSNextProto,
		ReadTimeout:  config.getReadTimeout(),
		WriteTimeout: config.getWriteTimeout(),
		Addr:         config.getServerAddr(),
		Handler:      engine,
	}
	impl := &serviceImpl{
		logger: logger,
		engine: engine,
		server: server,
	}
	return impl, nil
}

type serviceImpl struct {
	logger *zap.Logger
	engine *gin.Engine
	server *http.Server
}

func (impl *serviceImpl)NewRouterGroup(prefix string)*gin.RouterGroup{
	return impl.engine.Group(prefix)
}

func (impl *serviceImpl) ReleaseMode() {
	gin.SetMode(gin.ReleaseMode)
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
	impl.server.RegisterOnShutdown(onShutdown)
	go func() {
		err := impl.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			impl.logger.Error("HTTP服务异常关闭", zap.Error(err))
		}
		errorSign <- err
	}()
	return errorSign
}
