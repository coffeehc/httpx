// transport
package web

import (
	"github.com/coffeehc/logger"
	"github.com/coffeehc/render"
	"io"
	"net/http"
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

func transport_stream(r *render.Render, request *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	if reader, ok := v.(io.Reader); ok {
		return r.Stream(w, status, reader)
	}
	return Error_TransportType
}

func transport_binary(r *render.Render, request *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	if data, ok := v.([]byte); ok {
		return r.Data(w, status, data)
	}
	return Error_TransportType

}

func transport_html(r *render.Render, request *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	transportHtml, ok := v.(TransportHtml)
	if !ok {
		logger.Error("数据不是TransportHtml类型")
		return Error_TransportType
	}
	return r.HTML(w, status, transportHtml.templateName, transportHtml.binding, transportHtml.htmlOptiions...)
}

func transport_json(r *render.Render, request *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	return r.JSON(w, status, v)
}

func transport_jsonp(r *render.Render, request *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	callback := request.FormValue("callback")
	return r.JSONP(w, status, callback, v)
}

func transport_text(r *render.Render, request *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	if s, ok := v.(string); ok {
		return r.Text(w, status, s)
	}
	return Error_TransportType
}

func transport_xml(r *render.Render, request *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	return r.XML(w, status, v)
}

func transport_extend(r *render.Render, request *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	return nil
}

type TransportHtml struct {
	templateName string
	binding      interface{}
	htmlOptiions []render.HTMLOptions
}

func NewTransportHtml(templateName string, binding interface{}, htmlOptiions ...render.HTMLOptions) TransportHtml {
	return TransportHtml{templateName, binding, htmlOptiions}
}

//var (
//	Default_JsonTransport   = &JsonTransport{}
//	Default_StringTransport = StringTransport{}
//)
//
//type Transport interface {
//	Out(wirter io.Writer, data interface{}) error
//
//	ContentType() string
//}
//
//type StringTransport struct {
//}
//
//func (this StringTransport) Out(wirter io.Writer, data interface{}) error {
//	if str, ok := data.(string); ok {
//		wirter.Write([]byte(str))
//		return nil
//	}
//	_, err := wirter.Write([]byte(fmt.Sprintf("%#v", data)))
//	return err
//}
//
//func (this StringTransport) ContentType() string {
//	return "text/html;charset=UTF-8"
//}
//
//type JsonTransport struct {
//	jsonHandler *codec.JsonHandle
//	TimeFormat  codec.Ext
//}
//
//func (this *JsonTransport) Out(wirter io.Writer, data interface{}) error {
//	if this.jsonHandler == nil {
//		this.jsonHandler = new(codec.JsonHandle)
//		if this.TimeFormat == nil {
//			this.TimeFormat = TimeToStringConvert{"2006-01-02T15:04:05.999Z07:00"}
//		}
//		this.jsonHandler.SetExt(reflect.TypeOf(time.Time{}), 1, this.TimeFormat)
//	}
//	encode := codec.NewEncoder(wirter, this.jsonHandler)
//	err := encode.Encode(data)
//	if err != nil {
//		return err
//	}
//	return err
//}
//
//func (this *JsonTransport) ContentType() string {
//	return "json/application;charset=UTF-8"
//}
