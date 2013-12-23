package web

import "net/http"
import "fmt"

type Reply struct {
	request     *Request
	data        interface{}
	contentType string
	header      map[string]string
	tracsport   Tracsport
	statusCode  int
	redirect    string
	template    string
}

func NewReply(request *Request) *Reply {
	reply := new(Reply)
	reply.request = request
	reply.header = make(map[string]string)
	reply.statusCode = 200
	reply.tracsport = &Text
	return reply
}

func (this *Reply) With(data interface{}) *Reply {
	this.data = data
	return this
}

func (this *Reply) sendStatusCode(code int) *Reply {
	this.statusCode = code
	return this
}

func (this *Reply) Ok() *Reply {
	this.statusCode = 200
	return this
}

/**
 * Perform a 301 redirect (moved permanently) to the given uri.
 */
func (this *Reply) SeeOther(uri string) *Reply {
	this.statusCode = 301
	this.redirect = uri
	return this
}

/**
 * Perform a 302 redirect (moved permanently) to the given uri.
 */
func (this *Reply) Redirect(uri string) *Reply {
	this.statusCode = 302
	this.redirect = uri
	return this
}

/*
* 返回204没有内容
 */
func (this *Reply) noContent() *Reply {
	this.statusCode = 204
	return this
}

/*
 * 生成一个找不到页面的Replay
 */
func (this *Reply) NoFindPage() *Reply {
	this.statusCode = 404
	return this
}

func (this *Reply) forward(uri string) *Reply {
	this.request.URL.Path = uri
	return dispatcher.dispatch(this.request)
}

func (this *Reply) Header(header map[string]string) *Reply {
	for k, h := range header {
		if k != "" {
			this.header[k] = h
		}
	}
	return this
}

func (this *Reply) As(tracsport Tracsport) *Reply {
	this.tracsport = tracsport
	return this
}

func (this *Reply) writeResponse(w http.ResponseWriter, req *http.Request) error {
	w.Header().Set("Content-Type", this.tracsport.ContentType())
	for k, v := range this.header {
		w.Header().Set(k, v)
	}
	w.WriteHeader(this.statusCode)
	code := this.statusCode
	switch {
	case code >= 200 && code < 300:
		err := this.tracsport.Out(w, this)
		if err != nil {
			return err
		}
		break
	case code >= 300 && code < 400:
		http.Redirect(w, req, this.redirect, code)
		break
	case code >= 400 && code < 500:
		w.Write([]byte(fmt.Sprintf("%d 错误", code)))
		break
	case code >= 500:
		w.Write([]byte(fmt.Sprintf("%d 错误", code)))
		break
	default:
		break
	}
	return nil
}
