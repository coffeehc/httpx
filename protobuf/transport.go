package protobuf

import (
	"net/http"

	"github.com/coffeehc/web"
	"github.com/golang/protobuf/proto"
	"github.com/unrolled/render"
)

var (
	Transport_PB web.Transport = transport_pb
)

func transport_pb(r *render.Render, request *http.Request, w http.ResponseWriter, status int, v interface{}) error {
	if message, ok := v.(proto.Message); ok {
		head := render.Head{
			ContentType: "application/x-protobuf",
			Status:      status,
		}
		engine := ProtoBuf{
			Head: head,
		}
		return engine.Render(w, message)
	}
	return web.Error_TransportType
}

type ProtoBuf struct {
	render.Head
}

func (this ProtoBuf) Render(w http.ResponseWriter, message proto.Message) error {
	c := w.Header().Get(render.ContentType)
	if c != "" && this.Head.ContentType == "" {
		this.Head.ContentType = c
	}
	this.Head.Write(w)
	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
