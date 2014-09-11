// ResponseWriter
package web

import "net/http"

type responseWriter struct {
	base http.ResponseWriter
}

func NewResponseWriter(w http.ResponseWriter) http.ResponseWriter {
	return &responseWriter{base: w}
}

func (this *responseWriter) Header() http.Header {
	return this.base.Header()
}

func (this *responseWriter) Write(data []byte) (int, error) {
	return this.base.Write(data)
}

func (this *responseWriter) WriteHeader(code int) {
	this.base.WriteHeader(code)
}
