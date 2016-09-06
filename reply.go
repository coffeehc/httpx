// reply
package web

import (
	"net/http"

	"fmt"
	"github.com/coffeehc/logger"
	"github.com/unrolled/render"
)

type Reply interface {
	GetStatusCode() int
	SetStatusCode(statusCode int) Reply
	SetCookie(cookie http.Cookie) Reply
	SetHeader(key, value string) Reply
	AddHeader(key, value string) Reply
	DelHeader(key string) Reply
	GetHeader(key string) string
	Header() http.Header
	Redirect(code int, url string) Reply

	With(data interface{}) Reply
	As(transport Transport) Reply
	GetResponseWriter() http.ResponseWriter
	AdapterHttpHandler(adapter bool)
	GetRequestContext() RequestContext
}

type httpReply struct {
	statusCode         int
	data               interface{}
	header             http.Header
	cookies            []http.Cookie
	transport          Transport
	responseWriter     http.ResponseWriter
	adapterHttpHandler bool
	requestContext     RequestContext
}

func newHttpReply(responseWriter http.ResponseWriter, config *HttpServerConfig) *httpReply {
	return &httpReply{
		statusCode:     200,
		transport:      config.GetDefaultTransport(),
		cookies:        make([]http.Cookie, 0),
		responseWriter: responseWriter,
		header:         responseWriter.Header(),
	}
}

func (this *httpReply) GetRequestContext() RequestContext {
	if this.requestContext == nil {
		this.requestContext = make(RequestContext)
	}
	return this.requestContext
}

func (this *httpReply) AdapterHttpHandler(adapter bool) {
	this.adapterHttpHandler = adapter
}

func (this *httpReply) GetStatusCode() int {
	return this.statusCode
}

func (this *httpReply) SetStatusCode(statusCode int) Reply {
	this.statusCode = statusCode
	return this
}

func (this *httpReply) SetCookie(cookie http.Cookie) Reply {
	this.cookies = append(this.cookies, cookie)
	return this
}

func (this *httpReply) SetHeader(key, value string) Reply {
	this.header.Set(key, value)
	return this
}
func (this *httpReply) AddHeader(key, value string) Reply {
	this.header.Add(key, value)
	return this
}
func (this *httpReply) DelHeader(key string) Reply {
	this.header.Del(key)
	return this
}

func (this *httpReply) GetHeader(key string) string {
	return this.header.Get(key)
}

func (this *httpReply) Header() http.Header {
	return this.header
}

func (this *httpReply) Redirect(code int, url string) Reply {
	this.responseWriter.Header().Set("Location", url)
	this.statusCode = code
	return this
}

func (this *httpReply) With(data interface{}) Reply {
	this.data = data
	return this
}

func (this *httpReply) As(transport Transport) Reply {
	if transport != nil {
		this.transport = transport
	}
	return this
}

func (this *httpReply) GetResponseWriter() http.ResponseWriter {
	return this.responseWriter
}

//Reply 最后的清理工作
func (this *httpReply) finishReply(request *http.Request, render *render.Render) {
	if this.adapterHttpHandler {
		return
	}
	this.writeWarpHeader()
	if this.data == nil {
		this.data = ""
	}
	err := this.transport(render, request, this.GetResponseWriter(), this.GetStatusCode(), this.data)
	if err != nil {
		logger.Error("render error %#v", err)
		Transport_Text(render, request, this.GetResponseWriter(), 500, fmt.Sprintf("render error :%s", err))
	}
}

func (this *httpReply) writeWarpHeader() {
	header := this.Header()
	for _, cookie := range this.cookies {
		header.Set("Set-Cookie", cookie.String())
	}
}
