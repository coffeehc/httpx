package httpx

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

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
	Render(data interface{}) (io.ReadCloser, error)
}

func renderReply(w http.ResponseWriter, render Render, data interface{}) (io.ReadCloser, error) {
	contentType := render.ContentType()
	if contentType == "" {
		contentType = "text/plain"
	}
	w.Header().Set("content-type", contentType)
	return render.Render(data)

}

func toReadCloser(data interface{}) io.ReadCloser {
	switch v := data.(type) {
	case string:
		return ioutil.NopCloser(strings.NewReader(v))
	case *string:
		return ioutil.NopCloser(strings.NewReader(*v))
	case []byte:
		return ioutil.NopCloser(bytes.NewReader(v))
	case io.ReadCloser:
		return v
	case io.Reader:
		return ioutil.NopCloser(v)
	default:
		logger.Error("无法解析的数据类型,%#v", data)
	}
	return ioutil.NopCloser(bytes.NewReader([]byte{}))
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
func (render RenderJSON) Render(data interface{}) (io.ReadCloser, error) {
	v, err := ffjson.Marshal(data)
	if err != nil {
		return nil, err
	}
	return toReadCloser(v), nil
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
func (render RenderText) Render(data interface{}) (io.ReadCloser, error) {
	return toReadCloser(data), nil
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
func (render RenderXML) Render(data interface{}) (io.ReadCloser, error) {
	v, err := xml.Marshal(data)
	if err != nil {
		return nil, err
	}
	return toReadCloser(v), nil
}
