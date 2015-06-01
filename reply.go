// reply
package web

import (
	"bufio"
	"io"
	"net/http"

	"github.com/coffeehc/logger"
)

type Reply struct {
	statusCode  int
	data        interface{}
	headers     map[string]string
	cookies     []http.Cookie
	redirect    bool
	transport   Transport
	contentType string
	openStream  bool
	w           http.ResponseWriter
}

func newReply(w http.ResponseWriter) *Reply {
	return &Reply{statusCode: 200, transport: StringTransport{}, headers: make(map[string]string, 0), cookies: make([]http.Cookie, 0), w: w}
}

func (this *Reply) OpenStream() *Stream {
	this.SetContentType("text/plain")
	this.writeHeader()
	this.openStream = true
	return &Stream{w: this.w}
}

func (this *Reply) GetStatusCode() int {
	return this.statusCode
}

func (this *Reply) SetContentType(contentType string) *Reply {
	this.contentType = contentType
	return this
}
func (this *Reply) SetCode(statusCode int) *Reply {
	this.statusCode = statusCode
	return this
}

func (this *Reply) SetCookie(cookie http.Cookie) *Reply {
	this.cookies = append(this.cookies, cookie)
	return this
}

func (this *Reply) SetHeader(key, value string) *Reply {
	this.headers[key] = value
	return this
}

func (this *Reply) Redirect(code int, url string) *Reply {
	this.redirect = true
	this.headers["Location"] = url
	this.statusCode = code
	return this
}

func (this *Reply) With(data interface{}) *Reply {
	this.data = data
	return this
}

func (this *Reply) As(transport Transport) *Reply {
	if transport != nil {
		this.transport = transport
	}
	return this
}

func (this *Reply) write() {
	this.writeHeader()
	if this.redirect {
		return
	}
	if this.data != nil {
		if reader, ok := this.data.(io.Reader); ok {
			_, err := bufio.NewReader(reader).WriteTo(this.w)
			if err != nil {
				logger.Error("数据输出出现错误:%s", err)
			}
			if closer, ok := this.data.(io.Closer); ok {
				closer.Close()
			}
		} else {
			this.transport.Out(this.w, this.data)
		}
	}
}

func (this *Reply) writeHeader() {
	header := this.w.Header()
	for key, value := range this.headers {
		header.Set(key, value)
	}
	for _, cookie := range this.cookies {
		header.Set("Set-Cookie", cookie.String())
	}
	if this.contentType != "" {
		header.Set("Content-Type", this.contentType)
	} else {
		header.Set("Content-Type", this.transport.ContentType())
	}
	this.w.WriteHeader(this.statusCode)
}
