package web

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/coffeehc/logger"
)

const (
	PATH_SEPARATOR string = "/"
)

type routingDispatcher struct {
	routes map[string][]*route
}

func newRoutingDispatcher() *routingDispatcher {
	r := new(routingDispatcher)
	r.routes = make(map[string][]*route)
	return r
}

func (this *routingDispatcher) dispatch(request *http.Request, reply *Reply) {
	uri := request.URL.Path
	r := this.get(uri, strings.ToUpper(request.Method))
	if r == nil {
		FileHandler(request, reply)
		return
	}
	r.doMethod(request, reply)
}

type route struct {
	uri     string
	matcher pathMatcher
	//handlers map[string]interface{}
	requestMethod string
	handler       interface{}
}

type PathMap struct {
	pathMap map[string]string
}

func NewPathMap() *PathMap {
	return &PathMap{pathMap: make(map[string]string)}
}

func (this *PathMap) Get(key string) string {
	return this.pathMap[key]
}

func (this *route) callAction(req *http.Request, reply *Reply, pathMap *PathMap) {
	reply.Injector.Binding(pathMap, (*PathMap)(nil), "")
	reply.Injector.Invoke(this.handler)
}

func (this *route) doMethod(req *http.Request, reply *Reply) {
	matches := this.matcher.findMatches(req.URL.Path)
	//不处理Form里面的值，由程序自己处理掉
	this.callAction(req, reply, matches)
}

func newRoute(uri string, requestMethod string, handler interface{}) *route {
	rt := new(route)
	rt.uri = uri
	rt.matcher = newPathMatcherChain(uri)
	rt.requestMethod = requestMethod
	rt.handler = handler
	return rt
}

type pathMatcher interface {
	matches(incoming string) bool
	name() string
	findMatches(incoming string) *PathMap
}
type pathMatcherChain struct {
	pathMatcher
	path []pathMatcher
}

func newPathMatcherChain(path string) *pathMatcherChain {
	pmc := new(pathMatcherChain)
	pmc.path = toMatchChain(path)
	return pmc
}

func (this *pathMatcherChain) name() string {
	return ""
}
func (this *pathMatcherChain) matches(incoming string) bool {
	return this.findMatches(incoming) != nil
}

func (this *pathMatcherChain) findMatches(incoming string) *PathMap {
	pieces := strings.Split(incoming, PATH_SEPARATOR)
	if len(this.path) > len(pieces) {
		return nil
	}
	matches := NewPathMap()
	for i, pathMatcher := range this.path {
		if i == len(pieces) {
			if pathMatcher.matches("") {
				return matches
			}
			return nil
		}
		piece := pieces[i]
		if !pathMatcher.matches(piece) {
			return nil
		}
		name := pathMatcher.name()
		if len(name) != 0 {
			matches.pathMap[name] = piece
		}
	}
	if len(this.path) == len(pieces) {
		return matches
	}
	return nil
}

func toMatchChain(path string) []pathMatcher {
	pieces := strings.Split(path, PATH_SEPARATOR)
	matchers := make([]pathMatcher, len(pieces))
	for i, piece := range pieces {
		if strings.HasPrefix(piece, ":") {
			logger.Debug("Path中有参数的,需要使用GreedyPathMatcher匹配%s", piece)
			matchers[i] = newGreedyPathMatcher(strings.TrimLeft(piece, ":"))
		} else {
			matchers[i] = newSimplePathMatcher(piece)
		}
	}
	return matchers //不可变
}

type simplePathMatcher struct {
	pathMatcher
	path string
}

func newSimplePathMatcher(path string) *simplePathMatcher {
	return &simplePathMatcher{path: path}
}

func (this *simplePathMatcher) matches(incoming string) bool {
	return this.path == incoming
}
func (this *simplePathMatcher) name() string {
	return ""
}
func (this *simplePathMatcher) findMatches(incoming string) *PathMap {
	return NewPathMap()
}

type greedyPathMatcher struct {
	pathMatcher
	variable string
}

func newGreedyPathMatcher(piece string) *greedyPathMatcher {
	return &greedyPathMatcher{variable: piece}
}
func (this *greedyPathMatcher) matches(incoming string) bool {
	return true
}
func (this *greedyPathMatcher) name() string {
	return this.variable
}
func (this *greedyPathMatcher) findMatches(incoming string) *PathMap {
	return NewPathMap()
}
func checkFuncInterface(handler interface{}) {
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		panic("注册请求处理方法不是一个函数，不能注册")
	}
}
func (this *routingDispatcher) at(uri string, requestMethod string, handler interface{}) {
	checkFuncInterface(handler)
	key := firstPathElement(uri)
	rt := newRoute(uri, requestMethod, handler)
	if strings.HasPrefix(key, ":") {
		panic("不能使用第一路径随意匹配,如[/:xxx/xxx]")
	} else {
		rts := this.routes[key]
		if rts == nil {
			rts = make([]*route, 0)
		}
		rts = append(rts, rt)
		//这里没有排序啊,也没有做冲突处理，谁先注册就是谁先找到
		this.routes[key] = rts
	}
}

func (this *routingDispatcher) get(uri string, requestMethod string) *route {
	key := firstPathElement(uri)
	r := this.routes[key]
	if r != nil {
		for _, rt := range r {
			if rt.requestMethod == requestMethod && rt.matcher.matches(uri) {
				return rt
			}
		}
	}
	return nil
}

func firstPathElement(uri string) string {
	u := []rune(uri)
	index := strings.Index(string(u[1:]), "/")
	if index >= 0 {
		return string(u[1 : index+1])
	} else {
		return string(u[1:])
	}
}
