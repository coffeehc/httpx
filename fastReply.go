package web

import (
	"context"
	"net"
	"strings"

	"github.com/valyala/fasthttp"
	"unsafe"
)

func newReply(ctx *fasthttp.RequestCtx, _context context.Context, defaultTransportd Render) Reply {
	return &fastReply{ctx: ctx, context: _context, render: defaultTransportd}
}

type fastReply struct {
	ctx          *fasthttp.RequestCtx
	context      context.Context
	method       HttpMethod
	pathFragment PathFragment
	data         interface{}
	render       Render
}

func (this *fastReply) GetHttpMethod() HttpMethod {
	if this.method == "" {
		this.method = HttpMethod(strings.ToUpper(string(this.ctx.Request.Header.Method())))
	}
	return this.method
}

func (this *fastReply) GetFullURL() string {
	return string(this.ctx.Request.URI().FullURI())
}

func (this *fastReply) GetPath() string {
	return string(this.ctx.Request.URI().Path())
}

func (this *fastReply) GetRemoteAddr() net.Addr {
	return this.ctx.RemoteAddr()
}

func (this *fastReply) GetPathFragment() PathFragment {
	if this.pathFragment == nil {
		this.pathFragment = make(PathFragment, 0)
	}
	return this.pathFragment
}

func (this *fastReply) GetStatusCode() int {
	return this.ctx.Response.StatusCode()
}

func (this *fastReply) PutPathFragment(key, value string) Reply {
	if this.pathFragment == nil {
		this.pathFragment = make(PathFragment, 0)
	}
	this.pathFragment[key] = RequestParam(value)
	return this
}

func (this *fastReply) SetStatusCode(statusCode int) Reply {
	this.ctx.Response.SetStatusCode(statusCode)
	return this
}

func (this *fastReply) Redirect(code int, url string) Reply {
	this.ctx.Redirect(url, code)
	return this
}
func (this *fastReply) With(data interface{}) Reply {
	this.data = data
	return this
}
func (this *fastReply) As(render Render) Reply {
	this.render = render
	return this
}

func (this *fastReply) GetContext() context.Context {
	return this.context
}

func (this *fastReply) FinishReply() error {
	return this.render.Render(&(this.ctx.Response), this.data)
}

func (this *fastReply) GetRequestContext() *fasthttp.RequestCtx {
	return this.ctx
}

func (this *fastReply) GetQueryParam(key string) RequestParam {
	b := this.ctx.QueryArgs().Peek(key)
	if len(b) == 0 {
		return EmptyParam
	}
	return unsafeToRequestParam(b)
}

func (this *fastReply) GetPostParam(key string) RequestParam {
	b := this.ctx.Request.PostArgs().Peek(key)
	if len(b) == 0 {
		b = this.ctx.QueryArgs().Peek(key)
		if len(b) == 0 {
			return EmptyParam
		}
	}
	return unsafeToRequestParam(b)
}

func unsafeToRequestParam(b []byte) RequestParam {
	return *(*RequestParam)(unsafe.Pointer(&b))
}
