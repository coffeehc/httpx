package protobuf

import (
	"errors"
	"github.com/coffeehc/web"
	"github.com/golang/protobuf/proto"
	"net/http"
)

var (
	Render_PB web.Render = render_pb
)

func render_pb(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("content-type", "application/x-protobuf")
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
