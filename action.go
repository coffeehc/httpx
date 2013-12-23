// routes.go
package web

import "net/http"

const (
	METHOD_GET      string = "GET"
	METHOD_POST     string = "POST"
	METHOD_PUT      string = "PUT"
	METHOD_HEAD     string = "HEAD"
	METHOD_DELETE   string = "DELETE"
	METHOD_OPTIONSE string = "OPTIONSE"
	METHOD_TRACE    string = "TRACE"
)

type Request struct {
	*http.Request
	wenConf *WebConfig
}

//type ActionRegediter interface {
//	GetAction() *Action
//}

type Action struct {
	At      string
	methods []*method
}

type methodFunc func(req *Request, pathMap map[string]string) *Reply

type method struct {
	methodHandle methodFunc
	subAt        string
	httpMethod   string
}

func NewAction(at string) *Action {
	action := new(Action)
	action.At = at
	action.methods = make([]*method, 0)
	return action
}

func (this *Action) AddMethod(subAt string, httpMethod string, handle func(req *Request, pathMap map[string]string) *Reply) *Action {
	m := new(method)
	m.subAt = subAt
	m.httpMethod = httpMethod
	m.methodHandle = handle
	this.methods = append(this.methods, m)
	return this
}

func RegeditAction(action *Action) {
	dispatcher.serviceAt(action)
}
