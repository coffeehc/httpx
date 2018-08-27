package httpx

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type Config struct {
	ServerAddr     string
	ReadTimeout    int64       // 读的最大Timeout时间
	WriteTimeout   int64       // 写的最大Timeout时间
	MaxHeaderBytes int         // 请求头的最大长度
	TLSConfig      *tls.Config // 配置TLS
	TLSNextProto   map[string]func(*http.Server, *tls.Conn, http.Handler)
	ConnState      func(net.Conn, http.ConnState)
}

func (this *Config) getServerAddr() string {
	if this.ServerAddr == "" {
		this.ServerAddr = "0.0.0.0:8888"
	}
	return this.ServerAddr
}

func (this *Config) getReadTimeout() time.Duration {
	if this.ReadTimeout < 0 {
		this.ReadTimeout = 0
	}
	return time.Duration(this.ReadTimeout) * time.Second
}

func (this *Config) getWriteTimeout() time.Duration {
	if this.WriteTimeout < 0 {
		this.WriteTimeout = 0
	}
	return time.Duration(this.WriteTimeout) * time.Second
}

