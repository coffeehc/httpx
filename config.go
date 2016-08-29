package web

import (
	"time"

	"github.com/unrolled/render"
)

type HttpOption struct {
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

func (this *HttpOption) GetConcurrency() int {
	if this.Concurrency < 10000 {
		this.Concurrency = 100000
	}
	return this.Concurrency
}

func (this *HttpOption) GetDisableKeepalive() bool {
	return this.DisableKeepalive
}

func (this *HttpOption) GetReadBufferSize() int {
	if this.ReadBufferSize == 0 {
		this.ReadBufferSize = 1024 * 1024
	}
	return this.ReadBufferSize
}

func (this *HttpOption) GetWriteBufferSize() int {
	if this.WriteBufferSize == 0 {
		this.WriteBufferSize = 1024 * 1024
	}
	return this.WriteBufferSize
}

func (this *HttpOption) GetReadTimeout() time.Duration {
	if this._ReadTimeout == 0 {
		if this.ReadTimeout == 0 {
			this.ReadTimeout = 30
		}
		this._ReadTimeout = time.Second * time.Duration(this.ReadTimeout)
	}
	return this._ReadTimeout
}

func (this *HttpOption) GetWriteTimeout() time.Duration {
	if this._WriteTimeout == 0 {
		if this.WriteTimeout == 0 {
			this.WriteTimeout = 60
		}
		this._WriteTimeout = time.Second * time.Duration(this.WriteTimeout)
	}
	return this._WriteTimeout
}

func (this *HttpOption) GetMaxConnsPerIP() int {
	if this.MaxConnsPerIP == 0 {
		this.MaxConnsPerIP = 1000
	}
	return this.MaxConnsPerIP
}

func (this *HttpOption) GetMaxRequestsPerConn() int {
	if this.MaxRequestsPerConn == 0 {
		this.MaxRequestsPerConn = 1000
	}
	return this.MaxRequestsPerConn
}

func (this *HttpOption) GetMaxRequestBodySize() int {
	if this.MaxRequestBodySize == 0 {
		this.MaxRequestBodySize = 8 * 1024 * 1024
	}
	return this.MaxRequestBodySize
}

func (this *HttpOption) GetReduceMemoryUsage() bool {
	return this.ReduceMemoryUsage
}

type ServerConfig struct {
	Name       string `yaml:"name"`
	httpOption *HttpOption
	ServerAddr string
	//ConnState        func(net.Conn, http.ConnState)
	//HttpErrorLogout  io.Writer
	OpenTLS          bool
	CertFile         string
	KeyFile          string
	DefaultTransport Transport
	Render           *render.Render
}

func (this *ServerConfig) GetDefaultTransport() Transport {
	if this.DefaultTransport == nil {
		this.DefaultTransport = Transport_Text
	}
	return this.DefaultTransport
}

func (this *ServerConfig) GetRender() *render.Render {
	if this.Render == nil {
		this.Render = render.New()
	}
	return this.Render
}

func (this *ServerConfig) GetServerAddr() string {
	if this.ServerAddr == "" {
		this.ServerAddr = ":8888"
	}
	return this.ServerAddr
}

func (this *ServerConfig) GetName() string {
	if this.Name == "" {
		this.Name = "coffee"
	}
	return this.Name
}
