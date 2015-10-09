// transport
package web

import (
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/ugorji/go/codec"
)

var (
	Default_JsonTransport   = &JsonTransport{}
	Default_StringTransport = StringTransport{}
)

type Transport interface {
	Out(wirter io.Writer, data interface{}) error

	ContentType() string
}

type StringTransport struct {
}

func (this StringTransport) Out(wirter io.Writer, data interface{}) error {
	if str, ok := data.(string); ok {
		wirter.Write([]byte(str))
		return nil
	}
	_, err := wirter.Write([]byte(fmt.Sprintf("%#v", data)))
	return err
}

func (this StringTransport) ContentType() string {
	return "text/html;charset=UTF-8"
}

type JsonTransport struct {
	jsonHandler *codec.JsonHandle
	TimeFormat  codec.Ext
}

func (this *JsonTransport) Out(wirter io.Writer, data interface{}) error {
	if str, ok := data.(string); ok {
		wirter.Write([]byte(str))
		return nil
	}
	if this.jsonHandler == nil {
		this.jsonHandler = new(codec.JsonHandle)
		if this.TimeFormat != nil {
			this.jsonHandler.SetExt(reflect.TypeOf(time.Time{}), 1, this.TimeFormat)
		}
	}
	encode := codec.NewEncoder(wirter, this.jsonHandler)
	err := encode.Encode(data)
	if err != nil {
		return err
	}
	return err
}

func (this *JsonTransport) ContentType() string {
	return "json/application;charset=UTF-8"
}
