// reply
package web

import (
	"net/http"

	"github.com/coffeehc/render"
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
	AdapterHttpHander(adapter bool)
}

type httpReply struct {
	statusCode        int
	data              interface{}
	header            http.Header
	cookies           []http.Cookie
	transport         Transport
	responseWriter    http.ResponseWriter
	adapterHttpHander bool
}

func newHttpReply(responseWriter http.ResponseWriter) *httpReply {
	return &httpReply{
		statusCode:     200,
		transport:      Transport_Text,
		cookies:        make([]http.Cookie, 0),
		responseWriter: responseWriter,
		header:         responseWriter.Header(),
	}
}

func (this *httpReply) AdapterHttpHander(adapter bool) {
	this.adapterHttpHander = adapter
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
	if this.adapterHttpHander {
		return
	}
	this.writeWarpHeader()
	if this.data != nil {
		this.transport(render, request, this.GetResponseWriter(), this.GetStatusCode(), this.data)
		return
	}
	Transport_Text(render, request, this.GetResponseWriter(), this.GetStatusCode(), this.data)
}

func (this *httpReply) writeWarpHeader() {
	header := this.Header()
	for _, cookie := range this.cookies {
		header.Set("Set-Cookie", cookie.String())
	}
}
