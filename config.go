package httpx

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type Config struct {
	ServerAddr     string                                                 `yaml:"server_addr" json:"server_addr"`
	ReadTimeout    int64                                                  `yaml:"read_timeout" json:"read_timeout"`         // 读的最大Timeout时间
	WriteTimeout   int64                                                  `yaml:"write_timeout" json:"write_timeout"`       // 写的最大Timeout时间
	MaxHeaderBytes int                                                    `yaml:"max_header_bytes" json:"max_header_bytes"` // 请求头的最大长度
	TLSConfig      *tls.Config                                            `json:"-"`                                        // 配置TLS
	TLSNextProto   map[string]func(*http.Server, *tls.Conn, http.Handler) `json:"-"`
	ConnState      func(net.Conn, http.ConnState)                         `json:"-"`
}

func (impl *Config) getServerAddr() string {
	if impl.ServerAddr == "" {
		impl.ServerAddr = "0.0.0.0:8888"
	}
	return impl.ServerAddr
}

func (impl *Config) getReadTimeout() time.Duration {
	if impl.ReadTimeout < 0 {
		impl.ReadTimeout = 0
	}
	return time.Duration(impl.ReadTimeout) * time.Second
}

func (impl *Config) getWriteTimeout() time.Duration {
	if impl.WriteTimeout < 0 {
		impl.WriteTimeout = 0
	}
	return time.Duration(impl.WriteTimeout) * time.Second
}
