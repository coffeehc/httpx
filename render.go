package httpx

import (
	"io"

	"encoding/xml"
	"net/http"

	"github.com/coffeehc/logger"
	"github.com/pquerna/ffjson/ffjson"
)

var (
	//DefaultCharset 默认的字符集
	DefaultCharset = "utf-8"

	//DefaultRenderJSON 默认的 Json 渲染器
	DefaultRenderJSON = RenderJSON{Charset: DefaultCharset}
	//DefaultRenderXML 默认的 Xml 渲染器
	DefaultRenderXML = RenderXML{Charset: DefaultCharset}
	//DefaultRenderText 默认的 Text 渲染器
	DefaultRenderText = RenderText{Charset: DefaultCharset}
)

//Render 渲染器接口
type Render interface {
	ContentType() string
	Write(w http.ResponseWriter, data interface{}) error
}

func renderReply(w http.ResponseWriter, render Render, data interface{}) error {
	contentType := render.ContentType()
	if contentType == "" {
		contentType = "text/plain"
	}
	w.Header().Set("content-type", contentType)
	return render.Write(w, data)

}

func writeBody(w io.Writer, data interface{}) error {
	var err error
	switch v := data.(type) {
	case string:
		_, err = w.Write([]byte(v))
	case *string:
		_, err = w.Write([]byte(*v))
	case []byte:
		_, err = w.Write(v)
	case io.Reader:
		_, err = io.Copy(w, v)
		defer func() {
			if closer, ok := v.(io.Closer); ok {
				closer.Close()
			}
		}()
	case nil:
		logger.Warn("不返回任何信息")
	}
	return err
}

//RenderJSON json格式的渲染器
type RenderJSON struct {
	Charset string
}

//ContentType implement Render func
func (render RenderJSON) ContentType() string {
	return "application/json; charset=" + render.Charset
}

//Write implement Render func
func (render RenderJSON) Write(w http.ResponseWriter, data interface{}) error {
	v, err := ffjson.Marshal(data)
	if err != nil {
		return err
	}
	return writeBody(w, v)
}

//RenderText text格式的渲染器
type RenderText struct {
	Charset string
}

//ContentType implement Render func
func (render RenderText) ContentType() string {
	return "text/plain; charset=" + render.Charset
}

//Write implement Render func
func (render RenderText) Write(w http.ResponseWriter, data interface{}) error {
	return writeBody(w, data)
}

//RenderXML xml格式的渲染器
type RenderXML struct {
	Charset string
}

//ContentType implement Render func
func (render RenderXML) ContentType() string {
	return "text/plain; charset=" + render.Charset
}

//Write implement Render func
func (render RenderXML) Write(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("content-type", "application/xml; charset="+render.Charset)
	v, err := xml.Marshal(data)
	if err != nil {
		return err
	}
	return writeBody(w, v)
}
