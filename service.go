package httpx

import (
	"crypto/tls"
	"fmt"
	"github.com/coffeehc/base/log"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"net/http"
)

type Service interface {
	Start(onShutdown func()) <-chan error
	StartWithCertificate(cert tls.Certificate, onShutdown func()) <-chan error
	Shutdown() error
	GetEngine() *fiber.App
	NewRouterGroup(prefix string) fiber.Router
	GetServerAddress() string
}

func NewService(config *Config) Service {
	log.Debug(fmt.Sprintf("[%s]HTTP服务器配置", config.AppName))
	if config == nil {
		config = GetDefaultConfig("", "test")
	}
	engine := fiber.New(fiber.Config{
		Prefork:               config.Prefork,
		CaseSensitive:         config.CaseSensitive,
		StrictRouting:         config.StrictRouting,
		DisableStartupMessage: config.DisableStartupMessage,
		BodyLimit:             config.getBodyLimit(),
		Concurrency:           config.Concurrency,
		ServerHeader:          config.ServerHeader,
		AppName:               config.AppName,
		ReadTimeout:           config.getReadTimeout(),
		WriteTimeout:          config.getWriteTimeout(),
		IdleTimeout:           config.getIdleTimeout(),
		ReadBufferSize:        config.ReadBufferSize,
		WriteBufferSize:       config.WriteBufferSize,

		//TLSConfig: config.TLSConfig
		//TLSNextProto map[string]func(*http.Server, *tls.Conn, http.Handler)
		//ConnState    func(net.Conn, http.ConnState)

		Immutable:         config.Immutable,
		UnescapePath:      config.UnescapePath,
		ETag:              config.ETag,
		PassLocalsToViews: config.PassLocalsToViews,

		CompressedFileSuffix:         config.CompressedFileSuffix,
		ProxyHeader:                  config.ProxyHeader,
		GETOnly:                      config.GETOnly,
		DisableKeepalive:             config.DisableKeepalive,
		DisableDefaultDate:           config.DisableDefaultDate,
		DisableDefaultContentType:    config.DisableDefaultContentType,
		DisableHeaderNormalizing:     config.DisableHeaderNormalizing,
		StreamRequestBody:            config.StreamRequestBody,
		DisablePreParseMultipartForm: config.DisablePreParseMultipartForm,
		ReduceMemoryUsage:            config.ReduceMemoryUsage,
		Network:                      config.Network,
		EnableTrustedProxyCheck:      config.EnableTrustedProxyCheck,
		TrustedProxies:               config.TrustedProxies,
		EnableIPValidation:           config.EnableIPValidation,
		EnablePrintRoutes:            config.EnablePrintRoutes,
		Views:                        config.Views,
		ViewsLayout:                  config.ViewsLayout,
		ErrorHandler:                 config.ErrorHandler,
	})
	l, err := Listen(config.getServerAddr())
	if err != nil {
		log.Error(fmt.Sprintf("[%s]创建HTTP服务失败", config.AppName), zap.Error(err))
		return nil
	}
	config.ServerAddr = l.Addr().String()
	l.Close()
	impl := &serviceImpl{
		name:   config.AppName,
		config: config,
		engine: engine,
	}
	return impl
}

type serviceImpl struct {
	name   string
	config *Config
	engine *fiber.App
}

func (impl *serviceImpl) NewRouterGroup(prefix string) fiber.Router {
	return impl.engine.Group(prefix)
}

func (impl *serviceImpl) GetServerAddress() string {
	return impl.config.getServerAddr()
}

func (impl *serviceImpl) Shutdown() error {
	err := impl.engine.Shutdown()
	if err != nil {
		log.Error("关闭HttpServer失败", zap.Error(err))
	}
	return err
}

func (impl *serviceImpl) GetEngine() *fiber.App {
	return impl.engine
}

func (impl *serviceImpl) Start(onShutdown func()) <-chan error {
	errorSign := make(chan error, 1)
	go func() {
		err := impl.engine.Listen(impl.config.getServerAddr())
		if err != nil && err != http.ErrServerClosed {
			log.Error(fmt.Sprintf("[%s]HTTP服务异常关闭", impl.name), zap.Error(err))
		}
		log.Debug(fmt.Sprintf("[%s]HTTP服务关闭", impl.name))
		errorSign <- err
	}()
	log.Debug(fmt.Sprintf("[%s]HTTP服务启动", impl.name), zap.String("address", impl.config.getServerAddr()))
	return errorSign
}

func (impl *serviceImpl) StartWithCertificate(cert tls.Certificate, onShutdown func()) <-chan error {
	errorSign := make(chan error, 1)
	go func() {
		err := impl.engine.ListenTLSWithCertificate(impl.config.getServerAddr(), cert)
		if err != nil && err != http.ErrServerClosed {
			log.Error(fmt.Sprintf("[%s]HTTP服务异常关闭", impl.name), zap.Error(err))
		}
		log.Debug(fmt.Sprintf("[%s]HTTP服务关闭", impl.name))
		errorSign <- err
	}()
	log.Debug(fmt.Sprintf("[%s]HTTP服务启动", impl.name), zap.String("address", impl.config.getServerAddr()))
	return errorSign
}
