// extend
package web

import (
	"net/http"
)

type Stream struct {
	w http.ResponseWriter
}

func (this *Stream) CloseNotify() <-chan bool {
	return this.w.(http.CloseNotifier).CloseNotify()
}

func (this *Stream) Write(data []byte) (int, error) {
	return this.w.Write(data)
}

func (this *Stream) Flush() {
	this.w.(http.Flusher).Flush()
}
