package protobuf

import (
	"errors"
	"github.com/coffeehc/web"
	"github.com/golang/protobuf/proto"
	"github.com/valyala/fasthttp"
)

var (
	Render_PB web.Render = render_pb
)

func render_pb(response *fasthttp.Response, data interface{}) error {
	response.Header.SetContentType("application/x-protobuf")
	if message, ok := data.(proto.Message); ok {
		data, err := proto.Marshal(message)
		if err != nil {
			return err
		}
		response.SetBody(data)
	}
	return errors.New("data type is not proto.Message")
}
