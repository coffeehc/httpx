// transport
package web

import (
	"io"
	"net/http"

	"github.com/coffeehc/logger"
	"github.com/unrolled/render"
)

var (
	Transport_Stream Transport = transport_stream
	Transport_Binary Transport = transport_binary
	Transport_Html   Transport = transport_html
	Transport_Json   Transport = transport_json
	Transport_Jsonp  Transport = transport_jsonp
	Transport_Text   Transport = transport_text
	Transport_Xml    Transport = transport_xml
	Transport_Extend Transport = transport_extend
)

var Error_TransportType = HTTPERR_500("transport's interface type error")

type Transport func(r *render.Render, request *http.Request, w http.ResponseWriter, status int, v interface{}) error

func transport_stream(r *render.Render, rquest *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	if reader, ok := v.(io.Reader); ok {
		head := render.Head{
			ContentType: render.ContentBinary,
			Status:      status,
		}
		engine := Stream{
			Head: head,
		}
		return engine.Render(w, reader)
	}
	return Error_TransportType
}

func transport_binary(r *render.Render, rquest *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	if data, ok := v.([]byte); ok {
		return r.Data(w, status, data)
	}
	return Error_TransportType

}

func transport_html(r *render.Render, rquest *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	transportHtml, ok := v.(TransportHtml)
	if !ok {
		logger.Error("数据不是TransportHtml类型")
		return Error_TransportType
	}
	return r.HTML(w, status, transportHtml.templateName, transportHtml.binding, transportHtml.htmlOptions...)
}

func transport_json(r *render.Render, rquest *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	return r.JSON(w, status, v)
}

func transport_jsonp(r *render.Render, request *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	callback := request.FormValue("callback")
	return r.JSONP(w, status, callback, v)
}

func transport_text(r *render.Render, rquest *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	if s, ok := v.(string); ok {
		return r.Text(w, status, s)
	}
	return Error_TransportType
}

func transport_xml(r *render.Render, rquest *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	return r.XML(w, status, v)
}

func transport_extend(r *render.Render, rquest *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	return nil
}

type TransportHtml struct {
	templateName string
	binding      interface{}
	htmlOptions  []render.HTMLOptions
}

func NewTransportHtml(templateName string, binding interface{}, htmlOptions ...render.HTMLOptions) TransportHtml {
	return TransportHtml{templateName, binding, htmlOptions}
}

type Stream struct {
	render.Head
}

func (this Stream) Render(w http.ResponseWriter, reader io.Reader) error {
	c := w.Header().Get(render.ContentType)
	if c != "" {
		this.Head.ContentType = c
	}
	this.Head.Write(w)
	_, err := io.Copy(w, reader)
	return err
}
