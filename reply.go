package web

import (
	"bytes"
	"fmt"
	"html/template"
	"inject"
	"io"
	"net/http"
	"strings"

	"github.com/coffeehc/logger"
)

type Reply struct {
	inject.Injector
	data       io.Reader
	header     http.Header
	tracsport  Tracsport
	statusCode int
	//template   string
	sendSize int64
	cookies  []*http.Cookie
}

func NewReply(w http.ResponseWriter, parent inject.Injector) *Reply {
	reply := new(Reply)
	reply.Injector = inject.NewChild(parent)
	reply.header = w.Header()
	reply.statusCode = 200
	reply.tracsport = Html
	reply.cookies = make([]*http.Cookie, 0)
	return reply
}

func (this *Reply) GetSession() Session {
	if session, ok := this.GetInterface(Bind_Key_Session).(Session); ok {
		return session
	} else {
		return nil
	}
}

func (this *Reply) AddCookie(cookie *http.Cookie) {
	this.cookies = append(this.cookies, cookie)
}

func (this *Reply) DelCookie(cookie *http.Cookie) {
	cookie.MaxAge = -1
	this.AddCookie(cookie)
}

func (this *Reply) getTemplate(name string) *template.Template {
	baseTemp, ok := this.GetInterface(Bind_Key_Template).(*template.Template)
	if ok {
		if isProjectModel := this.GetInterface(Bind_Key_IsProjectModel).(bool); isProjectModel {
			return baseTemp.Lookup(name)
		} else {
			baseTemp = template.New("temp")
			return loadTemplate(this.GetInterface(Bind_Key_TemplateDir).(string), name, baseTemp)
		}
	}
	return nil
}

func (this *Reply) WithTemplate(data interface{}, name string) *Reply {
	template := this.getTemplate(name)
	if template == nil {
		logger.Debug("没有找到对应模版文件:%s", name)
		return this
	}
	buf := new(bytes.Buffer)
	t, err := template.Clone()
	if err != nil {
		logger.Error("模版克隆失败:%s", err)
		return this
	}
	t.Execute(buf, data)
	this.sendSize = (int64)(buf.Len())
	this.data = buf
	return this
}

/*
	将reader的内容输出到Response
	sendSize表示读取Reader里面的长度
	如果sendSize == -1则表示一直读取到EOF
*/
func (this *Reply) WithReader(data io.Reader, sendSize int64) *Reply {
	this.data = data
	this.sendSize = sendSize
	return this
}

/*
	将字符串返回给Response
*/
func (this *Reply) WithString(data string) *Reply {
	this.data = strings.NewReader(data)
	this.sendSize = -1
	return this
}

/*
	将byte返回给Response
*/
func (this *Reply) WithBytes(data []byte) *Reply {
	this.data = bytes.NewReader(data)
	this.sendSize = -1
	return this
}

func (this *Reply) SetStatusCode(code int) *Reply {
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
	this.header.Set("Location", uri)
	return this
}

/**
 * Perform a 302 redirect (moved permanently) to the given uri.
 */
func (this *Reply) Redirect(uri string) *Reply {
	this.statusCode = 302
	this.header.Set("Location", uri)
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
func (this *Reply) NoFindPage(request *http.Request) *Reply {
	this.statusCode = 404
	this.WithString(fmt.Sprintf("%s没有找到", request.RequestURI))
	return this
}

func (this *Reply) Forward(request *http.Request, uri string) {
	request.URL.Path = uri
	dispatcher.dispatch(request, this)
}

func (this *Reply) Header() http.Header {
	return this.header
}

func (this *Reply) As(tracsport Tracsport) *Reply {
	this.tracsport = tracsport
	return this
}
func (this *Reply) Error(err string, code int) {
	this.SetStatusCode(code)
	this.WithString(err)
}

func (this *Reply) writeResponse(w http.ResponseWriter, req *http.Request) {
	if _, haveType := w.Header()["Content-Type"]; !haveType {
		w.Header().Set("Content-Type", this.tracsport.ContentType())
	}
	if len(this.cookies) > 0 {
		for _, cookie := range this.cookies {
			http.SetCookie(w, cookie)
		}
	}
	code := this.statusCode
	if code >= 300 && code < 400 {
		http.Redirect(w, req, this.Header().Get("Location"), code)
	} else {
		reader := this.data
		defer func() {
			if closer, ok := reader.(io.Closer); ok {
				closer.Close()
			}
		}()
		w.WriteHeader(code)
		if this.sendSize < 0 {
			io.Copy(w, reader)
		} else {
			io.CopyN(w, reader, this.sendSize)
		}
	}
}
