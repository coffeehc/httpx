package web

import (
	"encoding/json"
	"io"

	"encoding/xml"
	"github.com/coffeehc/logger"
	"github.com/valyala/fasthttp"
)

var (
	DefaultCharset = "utf-8"

	Render_Json = render_Json{Charset: DefaultCharset}.Render
	Render_Xml  = render_Xml{Charset: DefaultCharset}.Render
	Render_Text = render_Text{Charset: DefaultCharset}.Render
)

func render(response *fasthttp.Response, data interface{}) error {
	switch v := data.(type) {
	case string:
		response.SetBodyString(v)
	case *string:
		response.SetBodyString(*v)
	case []byte:
		response.SetBody(v)
	case io.Reader:
		response.SetBodyStream(v, -1)
	case nil:
		logger.Warn("不返回任何信息")
	}
	return nil
}

type Render func(response *fasthttp.Response, data interface{}) error

type render_Json struct {
	Charset string
}

func (this render_Json) Render(response *fasthttp.Response, data interface{}) error {
	response.Header.SetContentType("application/json; charset=" + this.Charset)
	v, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return render(response, v)
}

type render_Text struct {
	Charset string
}

func (this render_Text) Render(response *fasthttp.Response, data interface{}) error {
	response.Header.SetContentType("text/plain; charset=" + this.Charset)
	return render(response, data)
}

type render_Xml struct {
	Charset string
}

func (this render_Xml) Render(response *fasthttp.Response, data interface{}) error {
	response.Header.SetContentType("application/xml; charset=" + this.Charset)
	v, err := xml.Marshal(data)
	if err != nil {
		return err
	}
	return render(response, v)
}
