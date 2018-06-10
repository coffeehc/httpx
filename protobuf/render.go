package protobuf

import (
	"errors"
	"io"

	"bytes"

	"io/ioutil"

	"git.xiagaogao.com/coffee/httpx"
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
func (render RenderProtobuf) Render(data interface{}) (io.ReadCloser, error) {
	if message, ok := data.(proto.Message); ok {
		data, err := proto.Marshal(message)
		if err != nil {
			return nil, err
		}
		return ioutil.NopCloser(bytes.NewBuffer(data)), nil
	}
	return nil, errors.New("data type is not proto.Message")
}
