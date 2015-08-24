// route
package web

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/coffeehc/logger"
)

const (
	PATH_SEPARATOR  = "/"
	WILDCARD_PREFIX = "{"
	WILDCARD_SUFFIX = "}"
	_conversion     = "[A-Za-z0-9]*"
)

type routingDispatcher struct {
	matcher pathMatcher
	filters []*filterWarp
}

func newRoutingDispatcher() *routingDispatcher {
	route := &routingDispatcher{matcher: pathMatcher{actionHandlerMap: make(map[string]_actionHandler)}, filters: make([]*filterWarp, 0)}
	route.addFilter(newServletStyleUriPatternMatcher("/*"), func(request *http.Request, reply *Reply, chain FilterChain) {
		defer func() {
			if ok := recover(); ok != nil {
				reply.SetCode(500).With(fmt.Sprintf("500:%#s", ok))
			}
		}()
		chain(request, reply)
	})
	return route
}

func (this *routingDispatcher) nextFilter(request *http.Request, index int) *filterWarp {
	if len(this.filters) <= index {
		return nil
	}
	filter := this.filters[index]
	if filter.matcher.match(request.URL.Path) {
		return filter
	}
	return this.nextFilter(request, filter.index)
}
func (route *routingDispatcher) addFilter(matcher uriPatternMatcher, actionFilter ActionFilter) {
	route.filters = append(route.filters, newFilterWarp(matcher, len(route.filters)+1, actionFilter, route))
}

func (route *routingDispatcher) handle(request *http.Request, reply *Reply) {
	defer func() {
		if err := recover(); err != nil {
			reply.SetCode(500).With(err)
		}
	}()
	handler := route.matcher.getActionHandler(request.URL.Path, strings.ToUpper(request.Method))
	if handler == nil {
		reply.SetCode(404).With("404:you are lost")
		logger.Error("Not found Handler for[%s] [%s]", strings.ToUpper(request.Method), request.URL.Path)
		return
	}
	handler.doAction(request, reply)
}

type actionHandler struct {
	method           string
	defineUri        string
	hasPathFragments bool
	uriConversions   map[string]int
	pathSize         int
	service          func(request *http.Request, PathFragments map[string]string, reply *Reply)
	exp              *regexp.Regexp
}

func (this *actionHandler) doAction(request *http.Request, reply *Reply) {
	requestUri := request.URL.Path
	param := make(map[string]string, 0)
	if this.hasPathFragments {
		paths := strings.Split(requestUri, PATH_SEPARATOR)
		if this.pathSize != len(paths) {
			panic(errors.New(logger.Error("需要解析的uri[%s]不匹配定义的uri[%s]", requestUri, this.defineUri)))
		}
		for name, index := range this.uriConversions {
			param[name] = paths[index]
		}
	}
	this.service(request, param, reply)
}

func (this *actionHandler) match(uri string) bool {
	return this.exp.MatchString(uri)
}

func buildActionHandler(action *defauleAction) (*actionHandler, error) {
	path := action.GetPath()
	if !strings.HasPrefix(path, PATH_SEPARATOR) {
		return nil, errors.New(logger.Error("定义的Uri必须是%s前缀", PATH_SEPARATOR))
	}
	paths := strings.Split(path, PATH_SEPARATOR)
	uriConversions := make(map[string]int, 0)
	conversionUri := ""
	pathSize := 0
	for index, p := range paths {
		pathSize++
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, WILDCARD_PREFIX) && strings.HasSuffix(p, WILDCARD_SUFFIX) {
			name := string([]byte(p)[len(WILDCARD_PREFIX) : len(p)-len(WILDCARD_SUFFIX)])
			uriConversions[name] = index
			p = _conversion
		}
		conversionUri += (PATH_SEPARATOR + p)
	}
	if conversionUri == "" {
		conversionUri = PATH_SEPARATOR
	}
	exp, err := regexp.Compile(conversionUri)
	if err != nil {
		return nil, err
	}
	return &actionHandler{exp: exp, method: action.GetMethod(), defineUri: action.GetPath(), hasPathFragments: len(uriConversions) > 0, uriConversions: uriConversions, pathSize: pathSize, service: action.Service}, nil
}

type _actionHandler []*actionHandler

func (this _actionHandler) Len() int {
	return len(this)
}
func (this _actionHandler) Less(i, j int) bool {
	uri1s := strings.Split(this[i].exp.String(), PATH_SEPARATOR)
	uri2s := strings.Split(this[j].exp.String(), PATH_SEPARATOR)
	if len(uri1s) != len(uri2s) {
		return len(uri1s) > len(uri2s)
	}
	for i, path1 := range uri1s {
		path2 := uri2s[i]
		if path1 != path2 {
			if path2 == _conversion {
				return false
			}
			return len(path1) < len(path2)
		}
	}
	return true
}
func (this _actionHandler) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

type pathMatcher struct {
	actionHandlerMap map[string]_actionHandler
}

func (this *pathMatcher) regeditAction(action *defauleAction) error {
	newActionHandler, err := buildActionHandler(action)
	if err != nil {
		logger.Error("添加Action出错:%s", err)
		return err
	}
	actionHandlers, ok := this.actionHandlerMap[action.GetMethod()]
	if !ok {
		actionHandlers = make(_actionHandler, 0)
	}
	for _, handler := range actionHandlers {
		if handler.exp.String() == newActionHandler.exp.String() {
			return errors.New(logger.Error("定义的uri[%s]与[%s]产生冲突,不能添加", handler.defineUri, newActionHandler.defineUri))
		}
	}
	this.actionHandlerMap[action.GetMethod()] = append(actionHandlers, newActionHandler)
	return nil
}

func (this *pathMatcher) getActionHandler(uri, method string) *actionHandler {
	actionHandlers, ok := this.actionHandlerMap[method]
	if !ok {
		return nil
	}
	paths := strings.Split(uri, PATH_SEPARATOR)
	pathSize := len(paths)
	for _, handler := range actionHandlers {
		if handler.pathSize == pathSize && handler.match(uri) {
			return handler
		}
	}
	return nil
}

func (this *pathMatcher) sort() {
	for _, actionHandlers := range this.actionHandlerMap {
		sort.Sort(actionHandlers)
	}
}
