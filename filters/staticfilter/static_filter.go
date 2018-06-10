package staticfilter

import (
	"net/http"

	"git.xiagaogao.com/coffee/httpx"
)

//RegisterStaticFilter register a static file handler to http server
func RegisterStaticFilter(server httpx.Server, uriPattern string, staticDir string) http.Handler {
	lastChar := uriPattern[len(uriPattern)-1]
	if lastChar != '*' {
		if lastChar != '/' {
			uriPattern += "/"
		}
		uriPattern = uriPattern + "*"
	}
	handler := http.StripPrefix(string(uriPattern[:len(uriPattern)-1]), http.FileServer(http.Dir(staticDir)))
	server.AddLastFilter(uriPattern, func(reply httpx.Reply, chain httpx.FilterChain) {
		reply.AdapterHTTPHandler(true)
		handler.ServeHTTP(reply.GetResponseWriter(), reply.GetRequest())
	})
	return handler
}
