package web

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"time"
)

type ServerConfig struct {
	ServerAddr        string
	ReadTimeout       time.Duration // 读的最大Timeout时间
	WriteTimeout      time.Duration // 写的最大Timeout时间
	MaxHeaderBytes    int           // 请求头的最大长度
	TLSConfig         *tls.Config   // 配置TLS
	TLSNextProto      map[string]func(*http.Server, *tls.Conn, http.Handler)
	ConnState         func(net.Conn, http.ConnState)
	HttpErrorLogout   io.Writer
	OpenTLS           bool
	CertFile          string
	KeyFile           string
	DefaultRender     Render
	DisabledKeepAlive bool
}

func (this *ServerConfig) getDisabledKeepAlive() bool {
	return this.DisabledKeepAlive
}

func (this *ServerConfig) getDefaultRender() Render {
	if this.DefaultRender == nil {
		this.DefaultRender = Default_Render_Text
	}
	return this.DefaultRender
}

func (this *ServerConfig) getServerAddr() string {
	if this.ServerAddr == "" {
		this.ServerAddr = "0.0.0.0:8888"
	}
	return this.ServerAddr
}

func (this *ServerConfig) getReadTimeout() time.Duration {
	if this.ReadTimeout < 0 {
		this.ReadTimeout = 0
	}
	return this.ReadTimeout
}

func (this *ServerConfig) getWriteTimeout() time.Duration {
	if this.WriteTimeout < 0 {
		this.WriteTimeout = 0
	}
	return this.WriteTimeout
}

func (this *ServerConfig) getMaxHeaderBytes() int {
	if this.MaxHeaderBytes < 0 {
		this.MaxHeaderBytes = 0
	}
	return this.MaxHeaderBytes
}
