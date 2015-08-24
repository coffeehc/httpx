// proxy
package proxy

import (
	"fmt"
	"net/http"

	"github.com/coffeehc/web"
)

type Proxy struct {
	host   string
	scheme string
	client *http.Client
}

func NewProxy(scheme, host string) *Proxy {
	return &Proxy{host, scheme, new(http.Client)}
}

func (this *Proxy) DoProxy(request *http.Request, reply *web.Reply, chain web.FilterChain) {
	request.URL.Scheme = this.scheme
	request.URL.Host = this.host
	request.Host = this.host
	request.RequestURI = ""
	resp, err := this.client.Do(request)
	if err != nil {
		reply.SetCode(500)
		reply.With(fmt.Sprintf("代理错误:%s", err))
		return
	}
	reply.SetCode(resp.StatusCode)
	for k, v := range resp.Header {
		reply.SetHeader(k, v[0])
	}
	reply.With(resp.Body)
}
