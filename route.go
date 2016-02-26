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

type routingDispatcher struct {
	matcher actionMatcher
	filter  *filterWarp
}

func newRoutingDispatcher() *routingDispatcher {
	route := &routingDispatcher{matcher: actionMatcher{actionHandlerMap: make(map[HttpMethod]actionHandlerList)}}
	route.filter = &filterWarp{
		matcher: newServletStyleUriPatternMatcher("/*"),
		actionFilter: func(request *http.Request, reply *Reply, chain FilterChain) {
			defer func() {
				if ok := recover(); ok != nil {
					reply.SetCode(500).With(fmt.Sprintf("500:%#s", ok))
				}
			}()
			chain(request, reply)
		},
		requestHandle: route.handle,
	}
	return route
}

func (route *routingDispatcher) addFilter(matcher uriPatternMatcher, actionFilter ActionFilter) {
	route.filter.addNextFilter(newFilterWarp(matcher, actionFilter))
}

func (route *routingDispatcher) handle(request *http.Request, reply *Reply) {
	handler := route.matcher.getActionHandler(request.URL.Path, HttpMethod(strings.ToUpper(request.Method)))
	if handler == nil {
		reply.SetCode(404).With("404:you are lost")
		logger.Error("Not found Handler for[%s] [%s]", strings.ToUpper(request.Method), request.URL.Path)
		return
	}
	handler.doAction(request, reply)
}
