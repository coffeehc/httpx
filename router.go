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
	_conversion     = "[^/]*"
)

type router struct {
	matcher       handlerMatcher
	filter        *filterWarp
	errorHandlers ErrorHandlers
}

func newRouter() *router {
	_router := &router{matcher: handlerMatcher{requestHandlerMap: make(map[HttpMethod]requestHandlerList)}}
	_router.errorHandlers = ErrorHandlers(make(map[int]RequestErrorHandler, 0))
	_router.filter = &filterWarp{
		matcher: newServletStyleUriPatternMatcher("/*"),
		requestFilter: func(request *http.Request, reply Reply, chain FilterChain) {
			defer func() {
				if err := recover(); err != nil {
					var httpErr *HttpError
					var ok bool
					if httpErr, ok = err.(*HttpError); !ok {
						httpErr = HTTPERR_500(fmt.Sprintf("%#s", err))
					}
					defer reply.SetStatusCode(httpErr.Code)
					if handler, ok := _router.errorHandlers[httpErr.Code]; ok {
						handler(request, httpErr, reply)
						return
					}
					reply.With(httpErr.Message).As(Transport_Json)
				}
			}()
			chain(request, reply)
		},
		filterChainFunc: _router.handle,
	}
	return _router
}

func (router *router) addFirstFilter(matcher UriPatternMatcher, actionFilter Filter) {
	oldFilter := router.filter
	newFilter := newFilterWarp(matcher, actionFilter)
	newFilter.nextFilter = oldFilter
	router.filter = newFilter

}

func (route *router) addLastFilter(matcher UriPatternMatcher, actionFilter Filter) {
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
