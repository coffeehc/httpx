package web

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"time"
)

type HttpServerConfig struct {
	ServerAddr        string
	ReadTimeout       time.Duration // 读的最大Timeout时间
	WriteTimeout      time.Duration // 写的最大Timeout时间
	MaxHeaderBytes    int           // 请求头的最大长度
	TLSConfig         *tls.Config   // 配置TLS
	TLSNextProto      map[string]func(*http.Server, *tls.Conn, http.Handler)
	ConnState         func(net.Conn, http.ConnState)
	HttpErrorLogout   io.Writer
	EnabledTLS        bool
	CertFile          string
	KeyFile           string
	DefaultRender     Render
	DisabledKeepAlive bool
}

func (this *HttpServerConfig) getDisabledKeepAlive() bool {
	return this.DisabledKeepAlive
}

func (this *HttpServerConfig) getDefaultRender() Render {
	if this.DefaultRender == nil {
		this.DefaultRender = Default_Render_Text
	}
	return this.DefaultRender
}

func (this *HttpServerConfig) getServerAddr() string {
	if this.ServerAddr == "" {
		this.ServerAddr = "0.0.0.0:8888"
	}
	return this.ServerAddr
}

func (this *HttpServerConfig) getReadTimeout() time.Duration {
	if this.ReadTimeout < 0 {
		this.ReadTimeout = 0
	}
	return this.ReadTimeout * time.Second
}

func (this *HttpServerConfig) getWriteTimeout() time.Duration {
	if this.WriteTimeout < 0 {
		this.WriteTimeout = 0
	}
	return this.WriteTimeout * time.Second
}

func (this *HttpServerConfig) getMaxHeaderBytes() int {
	if this.MaxHeaderBytes < 0 {
		this.MaxHeaderBytes = 0
	}
	return this.MaxHeaderBytes
}
