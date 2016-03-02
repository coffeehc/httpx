package static

import (
	"github.com/coffeehc/web"
	"net/http"
)

func RegeditStaticFilter(server *web.Server, uriPattern string, staticDir string) {
	lastChar := uriPattern[len(uriPattern)-1]
	if lastChar != '*' {
		if lastChar != '/' {
			uriPattern += "/"
		}
		uriPattern = uriPattern + "*"
	}
	handler := http.StripPrefix(string(uriPattern[:len(uriPattern)-1]), http.FileServer(http.Dir(staticDir)))
	server.AddLastFilter(uriPattern, func(request *http.Request, reply web.Reply, chain web.FilterChain) {
		reply.AdapterHttpHander(true)
		handler.ServeHTTP(reply.GetResponseWriter(), request)
	})
}
