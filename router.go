// route
package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/coffeehc/logger"
)

const (
	PATH_SEPARATOR  = "/"
	WILDCARD_PREFIX = "{"
	WILDCARD_SUFFIX = "}"
	_conversion     = "[A-Za-z0-9]*"
)

type router struct {
	matcher handlerMatcher
	filter  *filterWarp
}

func newRouter() *router {
	route := &router{matcher: handlerMatcher{requestHandlerMap: make(map[HttpMethod]requestHandlerList)}}
	route.filter = &filterWarp{
		matcher: newServletStyleUriPatternMatcher("/*"),
		requestFilter: func(request *http.Request, reply Reply, chain FilterChain) {
			defer func() {
				if ok := recover(); ok != nil {
					reply.SetStatusCode(500).With(fmt.Sprintf("500:%#s", ok))
				}
			}()
			chain(request, reply)
		},
		filterChainFunc: route.handle,
	}
	return route
}

func (route *router) addFilter(matcher UriPatternMatcher, actionFilter Filter) {
	route.filter.addNextFilter(newFilterWarp(matcher, actionFilter))
}

func (route *router) handle(request *http.Request, reply Reply) {
	handler := route.matcher.getActionHandler(request.URL.Path, HttpMethod(strings.ToUpper(request.Method)))
	if handler == nil {
		reply.SetStatusCode(404).With("404:you are lost")
		logger.Error("Not found Handler for[%s] [%s]", strings.ToUpper(request.Method), request.URL.Path)
		return
	}
	handler.doAction(request, reply)
}
