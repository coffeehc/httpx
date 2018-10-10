package httpx

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type Config struct {
	ServerAddr     string `yaml:"server_addr" json:"server_addr"`
	ReadTimeout    int64 `yaml:"read_timeout" json:"read_timeout"`    // 读的最大Timeout时间
	WriteTimeout   int64 `yaml:"write_timeout" json:"write_timeout"`    // 写的最大Timeout时间
	MaxHeaderBytes int `yaml:"max_header_bytes" json:"max_header_bytes"`     // 请求头的最大长度
	TLSConfig      *tls.Config `json:"-"`// 配置TLS
	TLSNextProto   map[string]func(*http.Server, *tls.Conn, http.Handler) `json:"-"`
	ConnState      func(net.Conn, http.ConnState) `json:"-"`
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
