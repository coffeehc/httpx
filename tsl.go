// extend
package web

import (
	"crypto/tls"
	"net/http"
)

func NewTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	config := &tls.Config{}
	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return config, nil
}

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
