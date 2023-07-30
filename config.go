package httpx

import (
	"crypto/tls"
	"github.com/gofiber/fiber/v2"
	"net"
	"net/http"
	"time"
)

func GetDefaultConfig(serverAddr string, appName string) *Config {
	return &Config{
		AppName:               appName,
		ServerAddr:            serverAddr,
		DisableStartupMessage: false,
		ServerHeader:          "Coffee",
	}
}

type Config struct {
	TLSConfig    *tls.Config
	TLSNextProto map[string]func(*http.Server, *tls.Conn, http.Handler)
	ConnState    func(net.Conn, http.ConnState)

	AppName        string `mapstructure:"app_name,omitempty" json:"app_name,omitempty"`
	ServerAddr     string `mapstructure:"server_addr,omitempty" json:"server_addr,omitempty"`
	ReadTimeoutMs  int64  `mapstructure:"read_timeout_ms,omitempty" json:"read_timeout_ms,omitempty"`
	WriteTimeoutMs int64  `mapstructure:"write_timeout_ms,omitempty" json:"write_timeout_ms,omitempty"`
	IdleTimeoutMs  int64  `mapstructure:"idle_timeout_ms,omitempty" json:"idle_timeout_ms,omitempty"`

	ReadBufferSize  int `mapstructure:"read_buffer_size,omitempty" json:"read_buffer_size,omitempty"`
	WriteBufferSize int `mapstructure:"write_buffer_size,omitempty" json:"write_buffer_size,omitempty"`
	MaxHeaderBytes  int `mapstructure:"max_header_bytes,omitempty" json:"max_header_bytes,omitempty"`

	Prefork           bool   `mapstructure:"prefork,omitempty" json:"prefork,omitempty"`
	ServerHeader      string `mapstructure:"server_header,omitempty" json:"server_header,omitempty"`
	StrictRouting     bool   `mapstructure:"strict_routing,omitempty" json:"strict_routing,omitempty"`
	CaseSensitive     bool   `mapstructure:"case_sensitive,omitempty" json:"case_sensitive,omitempty"`
	Immutable         bool   `mapstructure:"immutable,omitempty" json:"immutable,omitempty"`
	UnescapePath      bool   `mapstructure:"unescape_path,omitempty" json:"unescape_path,omitempty"`
	ETag              bool   `mapstructure:"e_tag,omitempty" json:"e_tag,omitempty"`
	BodyLimit         int    `mapstructure:"body_limit,omitempty" json:"body_limit,omitempty"`
	Concurrency       int    `mapstructure:"concurrency,omitempty" json:"concurrency,omitempty"`
	ViewsLayout       string `mapstructure:"views_layout,omitempty" json:"views_layout,omitempty"`
	PassLocalsToViews bool   `mapstructure:"pass_locals_to_views,omitempty" json:"pass_locals_to_views,omitempty"`

	CompressedFileSuffix         string   `mapstructure:"compressed_file_suffix,omitempty" json:"compressed_file_suffix,omitempty"`
	ProxyHeader                  string   `mapstructure:"proxy_header,omitempty" json:"proxy_header,omitempty"`
	GETOnly                      bool     `mapstructure:"get_only,omitempty" json:"get_only,omitempty"`
	DisableKeepalive             bool     `mapstructure:"disable_keepalive,omitempty" json:"disable_keepalive,omitempty"`
	DisableDefaultDate           bool     `mapstructure:"disable_default_date,omitempty" json:"disable_default_date,omitempty"`
	DisableDefaultContentType    bool     `mapstructure:"disable_default_content_type,omitempty" json:"disable_default_content_type,omitempty"`
	DisableHeaderNormalizing     bool     `mapstructure:"disable_header_normalizing,omitempty" json:"disable_header_normalizing,omitempty"`
	DisableStartupMessage        bool     `mapstructure:"disable_startup_message,omitempty" json:"disable_startup_message,omitempty"`
	StreamRequestBody            bool     `mapstructure:"stream_request_body,omitempty" json:"stream_request_body,omitempty"`
	DisablePreParseMultipartForm bool     `mapstructure:"disable_pre_parse_multipart_form,omitempty" json:"disable_pre_parse_multipart_form,omitempty"`
	ReduceMemoryUsage            bool     `mapstructure:"reduce_memory_usage,omitempty" json:"reduce_memory_usage,omitempty"`
	Network                      string   `mapstructure:"network,omitempty" json:"network,omitempty"`
	EnableTrustedProxyCheck      bool     `mapstructure:"enable_trusted_proxy_check,omitempty" json:"enable_trusted_proxy_check,omitempty"`
	TrustedProxies               []string `mapstructure:"trusted_proxies,omitempty" json:"trusted_proxies,omitempty"`
	EnableIPValidation           bool     `mapstructure:"enable_ip_validation,omitempty" json:"enable_ip_validation,omitempty"`
	EnablePrintRoutes            bool     `mapstructure:"enable_print_routes,omitempty" json:"enable_print_routes,omitempty"`

	Views fiber.Views `json:"-"`
}

func (impl *Config) getBodyLimit() int {
	if impl.BodyLimit == 0 {
		impl.BodyLimit = 1024 * 1024 * 16
	}
	return impl.BodyLimit
}

func (impl *Config) getServerAddr() string {
	if impl.ServerAddr == "" {
		impl.ServerAddr = "0.0.0.0:8888"
	}
	return impl.ServerAddr
}

func (impl *Config) getIdleTimeout() time.Duration {
	if impl.IdleTimeoutMs < 0 {
		impl.IdleTimeoutMs = 0
	}
	return time.Duration(impl.IdleTimeoutMs) * time.Millisecond
}

func (impl *Config) getReadTimeout() time.Duration {
	if impl.ReadTimeoutMs < 0 {
		impl.ReadTimeoutMs = 0
	}
	return time.Duration(impl.ReadTimeoutMs) * time.Millisecond
}

func (impl *Config) getWriteTimeout() time.Duration {
	if impl.WriteTimeoutMs < 0 {
		impl.WriteTimeoutMs = 0
	}
	return time.Duration(impl.WriteTimeoutMs) * time.Millisecond
}
