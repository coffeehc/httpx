// transport
package web

import (
	"encoding/json"
	"fmt"
	"io"
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
}

func (this JsonTransport) Out(wirter io.Writer, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = wirter.Write(bytes)
	return err
}

func (this JsonTransport) ContentType() string {
	return "json/application;charset=UTF-8"
}
