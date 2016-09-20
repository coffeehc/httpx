package protobuf

import (
	"errors"
	"github.com/coffeehc/web"
	"github.com/golang/protobuf/proto"
	"net/http"
)

var (
	Default_Render_PB web.Render = Render_protobuf{}
)

type Render_protobuf struct {
}

func (this Render_protobuf) ContentType() string {
	return "application/x-protobuf"
}

func (this Render_protobuf) Write(w http.ResponseWriter, data interface{}) error {
	if message, ok := data.(proto.Message); ok {
		data, err := proto.Marshal(message)
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		return err
	}
	return errors.New("data type is not proto.Message")
}
