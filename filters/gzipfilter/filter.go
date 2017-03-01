package gzipfilter

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/coffeehc/httpx"
)

//GZipFilter the gzip support
func GZipFilter(reply httpx.Reply, chain httpx.FilterChain) {
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
