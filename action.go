// routes.go
package web

import "path"

const (
	REQUEST_METHOD_GET      string = "GET"
	REQUEST_METHOD_POST     string = "POST"
	REQUEST_METHOD_PUT      string = "PUT"
	REQUEST_METHOD_HEAD     string = "HEAD"
	REQUEST_METHOD_DELETE   string = "DELETE"
	REQUEST_METHOD_OPTIONSE string = "OPTIONSE"
	REQUEST_METHOD_TRACE    string = "TRACE"
)

func HTTP_GET(uri string, handler interface{}) {
	dispatcher.at(uri, REQUEST_METHOD_GET, handler)
}
func HTTP_POST(uri string, handler interface{}) {
	dispatcher.at(uri, REQUEST_METHOD_POST, handler)
}
func HTTP_PUT(uri string, handler interface{}) {
	dispatcher.at(uri, REQUEST_METHOD_PUT, handler)
}
func HTTP_HEAD(uri string, handler interface{}) {
	dispatcher.at(uri, REQUEST_METHOD_HEAD, handler)
}
func HTTP_DELETE(uri string, handler interface{}) {
	dispatcher.at(uri, REQUEST_METHOD_DELETE, handler)
}
func HTTP_OPTIONSE(uri string, handler interface{}) {
	dispatcher.at(uri, REQUEST_METHOD_OPTIONSE, handler)
}
func HTTP_TRACE(uri string, handler interface{}) {
	dispatcher.at(uri, REQUEST_METHOD_TRACE, handler)
}

type HttpGroup struct {
	baseUri string
}

func NewHttpGroup(baseUri string) *HttpGroup {
	return &HttpGroup{baseUri: baseUri}
}

func (this *HttpGroup) GET(uri string, handler interface{}) {
	dispatcher.at(path.Join(this.baseUri, uri), REQUEST_METHOD_GET, handler)
}
func (this *HttpGroup) POST(uri string, handler interface{}) {
	dispatcher.at(path.Join(this.baseUri, uri), REQUEST_METHOD_POST, handler)
}
func (this *HttpGroup) PUT(uri string, handler interface{}) {
	dispatcher.at(path.Join(this.baseUri, uri), REQUEST_METHOD_PUT, handler)
}
func (this *HttpGroup) HEAD(uri string, handler interface{}) {
	dispatcher.at(path.Join(this.baseUri, uri), REQUEST_METHOD_HEAD, handler)
}
func (this *HttpGroup) DELETE(uri string, handler interface{}) {
	dispatcher.at(path.Join(this.baseUri, uri), REQUEST_METHOD_DELETE, handler)
}
func (this *HttpGroup) OPTIONSE(uri string, handler interface{}) {
	dispatcher.at(path.Join(this.baseUri, uri), REQUEST_METHOD_OPTIONSE, handler)
}
func (this *HttpGroup) TRACE(uri string, handler interface{}) {
	dispatcher.at(path.Join(this.baseUri, uri), REQUEST_METHOD_TRACE, handler)
}
