package protobuf

import (
	"errors"
	"net/http"

	"github.com/coffeehc/httpx"
	"github.com/golang/protobuf/proto"
)

var (
	//DefaultRenderPB the protobuf'srender instance
	DefaultRenderPB httpx.Render = RenderProtobuf{}
)

//RenderProtobuf httpx.Render implement by protobuf
type RenderProtobuf struct {
}

//ContentType implement httpx.Render interface
func (render RenderProtobuf) ContentType() string {
	return "application/x-protobuf"
}

//Write implement httpx.Render interface
func (render RenderProtobuf) Write(w http.ResponseWriter, data interface{}) error {
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
