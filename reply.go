// reply
package web

import (
	"net/http"

	"fmt"
	"github.com/coffeehc/logger"
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
	AddPathFragment(k, v string)

	With(data interface{}) Reply
	As(render Render) Reply

	GetRequest() *http.Request
	GetResponseWriter() http.ResponseWriter
	GetPathFragment() PathFragment
	AdapterHttpHandler(adapter bool)
}

type httpReply struct {
	statusCode         int
	data               interface{}
	header             http.Header
	cookies            []http.Cookie
	render             Render
	request            *http.Request
	responseWriter     http.ResponseWriter
	adapterHttpHandler bool
	pathFragment       PathFragment
}

func newHttpReply(request *http.Request, w http.ResponseWriter, config *ServerConfig) *httpReply {
	return &httpReply{
		statusCode:     200,
		render:         config.getDefaultRender(),
		cookies:        make([]http.Cookie, 0),
		request:        request,
		responseWriter: w,
		header:         w.Header(),
	}
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

func (this *httpReply) As(render Render) Reply {
	if render != nil {
		this.render = render
	}
	return this
}

func (this *httpReply) GetRequest() *http.Request {
	return this.request
}

func (this *httpReply) GetResponseWriter() http.ResponseWriter {
	return this.responseWriter
}

func (this *httpReply) GetPathFragment() PathFragment {
	return this.pathFragment
}

func (this *httpReply) AddPathFragment(k, v string) {
	if this.pathFragment == nil {
		this.pathFragment = make(PathFragment, 0)
	}
	this.pathFragment[k] = RequestParam(v)
}

//Reply 最后的清理工作
func (this *httpReply) finishReply() {
	if this.adapterHttpHandler {
		return
	}
	this.writeWarpHeader()
	if this.data == nil {
		this.data = ""
	}
	err := renderReply(this.GetResponseWriter(), this.render, this.data)
	if err != nil {
		this.SetStatusCode(500)
		logger.Error("render error %#v", err)
		renderReply(this.GetResponseWriter(), Default_Render_Text, fmt.Sprintf("render error :%s", err))
	}
}

func (this *httpReply) writeWarpHeader() {
	header := this.Header()
	for _, cookie := range this.cookies {
		header.Set("Set-Cookie", cookie.String())
	}
}
