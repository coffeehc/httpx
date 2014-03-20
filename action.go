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

type ActionHandler func(req *Request, pathMap map[string]string, reply *Reply)

type method struct {
	methodHandler ActionHandler
	subAt         string
	httpMethod    string
}

func NewAction(at string) *Action {
	action := new(Action)
	action.At = at
	action.methods = make([]*method, 0)
	return action
}

func (this *Action) AddMethod(subAt string, httpMethod string, handler ActionHandler) {
	m := new(method)
	m.subAt = subAt
	m.httpMethod = httpMethod
	m.methodHandler = handler
	this.methods = append(this.methods, m)
}

func RegeditAction(action *Action) {
	dispatcher.serviceAt(action)
}
