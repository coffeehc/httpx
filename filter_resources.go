// filter_static
package web

import (
	"net/http"
	"strings"
)

type ResourceFilter struct {
	Pattern   string
	Prefix    string
	prefixLen int
}

func (this *ResourceFilter) GetPattern() string {
	return this.Pattern
}

func (this *ResourceFilter) Init(conf *WebConfig) {
	this.prefixLen = len(this.Prefix)
}

func (this *ResourceFilter) DoFilter(req *http.Request, reply *Reply, chain FilterChain) {
	path := req.URL.Path
	if this.prefixLen != 0 && strings.HasPrefix(path, this.Prefix) {
		req.URL.Path = path[this.prefixLen:]
		FileHandler(req, reply)
	} else {
		chain.DoFilter(req, reply)
	}
}

func (this *ResourceFilter) Destroy() {
}
