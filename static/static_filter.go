package static

import (
	"net/http"

	"github.com/coffeehc/web"
)

func RegisterStaticFilter(server *web._Server, uriPattern string, staticDir string) http.Handler {
	lastChar := uriPattern[len(uriPattern)-1]
	if lastChar != '*' {
		if lastChar != '/' {
			uriPattern += "/"
		}
		uriPattern = uriPattern + "*"
	}
	handler := http.StripPrefix(string(uriPattern[:len(uriPattern)-1]), http.FileServer(http.Dir(staticDir)))
	server.AddLastFilter(uriPattern, func(request *http.Request, reply web.Reply, chain web.FilterChain) {
		reply.AdapterHttpHandler(true)
		handler.ServeHTTP(reply.GetResponseWriter(), request)
	})
	return handler
}
