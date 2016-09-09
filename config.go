package web

import (
	"time"

	"github.com/valyala/fasthttp"
)

type ServerOption struct {
	Concurrency        int   `yaml:"concurrency"`
	DisableKeepalive   bool  `yaml:"disable_keepalive"`
	ReadBufferSize     int   `yaml:"read_buffer_size"`
	WriteBufferSize    int   `yaml:"write_buffer_size"`
	ReadTimeout        int64 `yaml:"read_timeout"` //unit:second
	_ReadTimeout       time.Duration
	WriteTimeout       int64 `yaml:"write_timeout"` //unit:second
	_WriteTimeout      time.Duration
	MaxConnsPerIP      int  `yaml:"max_conns_per_ip"`
	MaxRequestsPerConn int  `yaml:"max_requests_per_conn"`
	MaxRequestBodySize int  `yaml:"max_request_body_size"`
	ReduceMemoryUsage  bool `yaml:"reduce_memory_usage"`
}

func (this *ServerOption) GetConcurrency() int {
	if this.Concurrency < 10000 {
		this.Concurrency = fasthttp.DefaultConcurrency //100000
	}
	return this.Concurrency
}

func (this *ServerOption) GetDisableKeepalive() bool {
	return this.DisableKeepalive
}

func (this *ServerOption) GetReadBufferSize() int {
	if this.ReadBufferSize == 0 {
		this.ReadBufferSize = 1024 * 1024
	}
	return this.ReadBufferSize
}

func (this *ServerOption) GetWriteBufferSize() int {
	if this.WriteBufferSize == 0 {
		this.WriteBufferSize = 1024 * 1024
	}
	return this.WriteBufferSize
}

func (this *ServerOption) GetReadTimeout() time.Duration {
	if this._ReadTimeout == 0 {
		if this.ReadTimeout == 0 {
			this.ReadTimeout = 30
		}
		this._ReadTimeout = time.Second * time.Duration(this.ReadTimeout)
	}
	return this._ReadTimeout
}

func (this *ServerOption) GetWriteTimeout() time.Duration {
	if this._WriteTimeout == 0 {
		if this.WriteTimeout == 0 {
			this.WriteTimeout = 60
		}
		this._WriteTimeout = time.Second * time.Duration(this.WriteTimeout)
	}
	return this._WriteTimeout
}

func (this *ServerOption) GetMaxConnsPerIP() int {
	if this.MaxConnsPerIP == 0 {
		this.MaxConnsPerIP = 1000
	}
	return this.MaxConnsPerIP
}

func (this *ServerOption) GetMaxRequestsPerConn() int {
	if this.MaxRequestsPerConn == 0 {
		this.MaxRequestsPerConn = 1000
	}
	return this.MaxRequestsPerConn
}

func (this *ServerOption) GetMaxRequestBodySize() int {
	if this.MaxRequestBodySize == 0 {
		this.MaxRequestBodySize = 8 * 1024 * 1024
	}
	return this.MaxRequestBodySize
}

func (this *ServerOption) GetReduceMemoryUsage() bool {
	return this.ReduceMemoryUsage
}

type HttpServerConfig struct {
	Name          string `yaml:"name"`
	ServerOption  *ServerOption
	ServerAddr    string
	OpenTLS       bool
	CertFile      string
	KeyFile       string
	DefaultRender Render
}

func (this *HttpServerConfig) GetServerOption() *ServerOption {
	if this.ServerOption == nil {
		this.ServerOption = &ServerOption{}
	}
	return this.ServerOption
}

func (this *HttpServerConfig) GetDefaultRender() Render {
	if this.DefaultRender == nil {
		this.DefaultRender = Render_Text
	}
	return this.DefaultRender
}

func (this *HttpServerConfig) GetServerAddr() string {
	if this.ServerAddr == "" {
		this.ServerAddr = ":8888"
	}
	return this.ServerAddr
}

func (this *HttpServerConfig) GetName() string {
	if this.Name == "" {
		this.Name = "coffee's httpServer"
	}
	return this.Name
}
