package gzipfilter

import (
	"compress/gzip"
	"github.com/coffeehc/web"
	"io"
	"net/http"
	"strings"
)

func GZipFilter(reply web.Reply, chain web.FilterChain) {
	w := reply.GetResponseWriter()
	r := reply.GetRequest()
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		reply.WarpResponseWriter(&gzipResponseWriter{Writer: gz, ResponseWriter: w})
	}
	chain(reply)
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
