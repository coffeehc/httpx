// reply
package web

import (
	"io"
	"net/http"

	"github.com/coffeehc/logger"
)

type Reply interface {
	GetStatusCode() int
	SetStatusCode(statusCode int) Reply
	GetContentType() string
	SetContentType(contentType string) Reply
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
	redirect          bool
	transport         Transport
	contentType       string
	openStream        bool
	responseWriter    http.ResponseWriter
	adapterHttpHander bool
}

func newHttpReply(responseWriter http.ResponseWriter) *httpReply {
	return &httpReply{
		statusCode:     200,
		transport:      Default_StringTransport,
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

func (this *httpReply) GetContentType() string {
	return this.contentType
}

func (this *httpReply) SetContentType(contentType string) Reply {
	this.contentType = contentType
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
	this.redirect = true
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
func (this *httpReply) finishReply() {
	if this.adapterHttpHander {
		return
	}
	this.writeWarpHeader()
	this.responseWriter.WriteHeader(this.GetStatusCode())
	if this.redirect {
		return
	}
	if this.data != nil {
		if reader, ok := this.data.(io.Reader); ok {
			_, err := io.Copy(this.responseWriter, reader)
			if err != nil {
				logger.Error("数据输出出现错误:%s", err)
				return
			}
			if limitReader, ok := this.data.(*io.LimitedReader); ok {
				reader = limitReader.R
				if closer, ok := reader.(io.Closer); ok {
					closer.Close()
				}
			}
			if closer, ok := this.data.(io.Closer); ok {
				closer.Close()
			}
		} else {
			err := this.transport.Out(this.responseWriter, this.data)
			if err != nil {
				logger.Error("数据序列化失败:%s", err)
			}
		}
	}
}

func (this *httpReply) writeWarpHeader() {
	header := this.Header()
	for _, cookie := range this.cookies {
		header.Set("Set-Cookie", cookie.String())
	}
	if this.contentType != "" {
		header.Set("Content-Type", this.contentType)
	} else {
		header.Set("Content-Type", this.transport.ContentType())
	}
}
