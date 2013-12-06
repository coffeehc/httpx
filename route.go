package web

import "strings"
import "logger"

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

func (this *routingDispatcher) dispatch(request *Request) *Reply {
	uri := request.URL.Path
	logger.Debugf("开始匹配URI:%s", uri)
	r := this.get(uri)
	if r == nil {
		return nil
	}
	return r.doMethod(request, r)
}

type route struct {
	uri       string
	matcher   pathMatcher
	headless  bool
	extension bool
	methods   map[string]methodFunc
}

func (this *route) callAction(req *Request, r *route, pathMap map[string]string) *Reply {
	tuple := this.methods[req.Method]
	var replay *Reply
	if tuple != nil {
		logger.Debug("开始执行真正的方法")
		//TODO解析URLValue
		replay = tuple(req, pathMap)
	} else {
		replay = NewReply(req).NoFindPage()
	}
	return replay
}

func (this *route) doMethod(req *Request, r *route) *Reply {
	matches := this.matcher.findMatches(req.URL.Path)
	return this.callAction(req, r, matches)
}

func (this *route) isHeadless() bool {
	return this.headless
}

func newRoute(uri string, matcher pathMatcher, headless, extension bool, action *Action) *route {
	logger.Debugf("创建一个新的RouteTuple:%s", uri)
	rt := new(route)
	rt.uri = uri
	rt.matcher = matcher
	rt.headless = headless
	rt.extension = extension
	methods := action.methods
	rt.methods = make(map[string]methodFunc, len(methods))
	for _, method := range methods {
		logger.Debugf("注册了[%s]的[%s]方法", uri, method.httpMethod)
		rt.methods[method.httpMethod] = method.methodHandle
	}
	return rt
}

type pathMatcher interface {
	matches(incoming string) bool
	name() string
	findMatches(incoming string) map[string]string
}
type pathMatcherChain struct {
	pathMatcher
	path []pathMatcher
}

func newPathMatcherChain(path string) *pathMatcherChain {
	pmc := new(pathMatcherChain)
	logger.Debugf("创建一个新的PathMatcherChain:%s", path)
	pmc.path = toMatchChain(path)
	return pmc
}

func (this *pathMatcherChain) name() string {
	return ""
}
func (this *pathMatcherChain) matches(incoming string) bool {
	return this.findMatches(incoming) != nil
}

func (this *pathMatcherChain) findMatches(incoming string) map[string]string {
	pieces := strings.Split(incoming, PATH_SEPARATOR)
	if len(this.path) > len(pieces) {
		return nil
	}
	matches := make(map[string]string)
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
			matches[name] = piece
		}
	}
	if len(this.path) == len(pieces) {
		return matches
	}
	return nil
}

func toMatchChain(path string) []pathMatcher {
	pieces := strings.Split(path, PATH_SEPARATOR)
	logger.Debugf("分解Path:%v", pieces)
	matchers := make([]pathMatcher, len(pieces))
	for i, piece := range pieces {
		if strings.HasPrefix(piece, ":") {
			logger.Debugf("Path中有参数的,需要使用GreedyPathMatcher匹配%s", piece)
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
func (this *simplePathMatcher) findMatches(incoming string) map[string]string {
	return make(map[string]string)
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
func (this *greedyPathMatcher) findMatches(incoming string) map[string]string {
	return make(map[string]string)
}

func (this *routingDispatcher) serviceAt(action *Action) {
	baseUri := action.At
	if len(baseUri) == 0 {
		panic("baseUri必须设置")
	}
	for _, m := range action.methods {
		subPath := m.subAt
		if len(subPath) != 0 && (!strings.HasPrefix(subPath, "/") || len(subPath) == 1) {
			panic("subPath 必须以\\开头且不能直接为'\\'")
		}
		this.at(baseUri+subPath, action, true)
	}
	this.at(baseUri, action, true)
}

func (this *routingDispatcher) at(uri string, action *Action, headless bool) {
	key := firstPathElement(uri)
	logger.Debugf("开始处理%s开头的路径At", key)
	rt := newRoute(uri, newPathMatcherChain(uri), headless, false, action)
	//需要加锁么
	if strings.HasPrefix(key, ":") {
		panic("不能使用第一路径随意匹配,如[/:xxx/xxx]")
	} else {
		rts := this.routes[key]
		if rts == nil {
			logger.Debugf("没有找到%s对应的routes数组,新建一个", key)
			rts = make([]*route, 0)
		}
		rts = append(rts, rt)
		//这里没有排序啊
		this.routes[key] = rts
	}
}

func (this *routingDispatcher) get(uri string) *route {
	key := firstPathElement(uri)
	logger.Debugf("获取了Path的第一个路径:%s", key)
	r := this.routes[key]
	logger.Debugf("找到了%d", len(r))
	if r != nil {
		for _, rt := range r {
			if rt.matcher.matches(uri) {
				logger.Debugf("找到了对应Route,%v", r)
				return rt
			}
		}
	}
	logger.Debugf("没有找到对应的route,返回空")
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
