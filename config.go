package httpx

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"time"
)

// Config http server Config
type Config struct {
	ServerAddr        string        `json:"server_addr" yaml:"server_addr"`           //ServerAddr server地址
	ReadTimeout       time.Duration `json:"read_timeout" yaml:"read_timeout"`         // 读的最大Timeout时间
	WriteTimeout      time.Duration `json:"write_timeout" yaml:"write_timeout"`       // 写的最大Timeout时间
	MaxHeaderBytes    int           `json:"max_header_bytes" yaml:"max_header_bytes"` // 请求头的最大长度
	TLSConfig         *tls.Config   // 配置TLS
	TLSNextProto      map[string]func(*http.Server, *tls.Conn, http.Handler)
	ConnState         func(net.Conn, http.ConnState)
	HTTPErrorLogout   io.Writer
	DefaultRender     Render
	KeepAliveDuration time.Duration `json:"keep_alive_duration" yaml:"keep_alive_duration"`
	IdleTimeout       time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
}

func (config *Config) getKeepAliveDuration() time.Duration {
	if config.KeepAliveDuration == 0 {
		config.KeepAliveDuration = 3 * time.Second
	}
	return config.KeepAliveDuration
}

func (config *Config) getDefaultRender() Render {
	if config.DefaultRender == nil {
		config.DefaultRender = DefaultRenderText
	}
	return config.DefaultRender
}

func (config *Config) getServerAddr() string {
	if config.ServerAddr == "" {
		config.ServerAddr = "0.0.0.0:8888"
	}
	return config.ServerAddr
}

func (config *Config) getReadTimeout() time.Duration {
	if config.ReadTimeout < 0 {
		config.ReadTimeout = 0
	}
	return config.ReadTimeout * time.Second
}

func (config *Config) getWriteTimeout() time.Duration {
	if config.WriteTimeout < 0 {
		config.WriteTimeout = 0
	}
	return config.WriteTimeout * time.Second
}

func (config *Config) getMaxHeaderBytes() int {
	if config.MaxHeaderBytes < 0 {
		// TODO
		config.MaxHeaderBytes = 0
	}
	return config.MaxHeaderBytes
}

func (config *Config) getIdleTimeout() time.Duration {
	if config.IdleTimeout < 0 {
		config.IdleTimeout = 0
	}
	return config.IdleTimeout * time.Second
}
