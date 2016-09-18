// route
package web

import (
	"fmt"

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
		requestFilter: func(reply Reply, chain FilterChain) {
			defer func() {
				if err := recover(); err != nil {
					var httpErr *HttpError
					var ok bool
					if httpErr, ok = err.(*HttpError); !ok {
						httpErr = HTTPERR_500(fmt.Sprintf("%#s", err))
					}
					reply.SetStatusCode(httpErr.Code)
					if handler, ok := _router.errorHandlers[httpErr.Code]; ok {
						handler(httpErr, reply)
						return
					}
					reply.With(httpErr.Message).As(Render_Json)
				}
			}()
			chain(reply)
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

func (route *router) handle(reply Reply) {
	path := reply.GetPath()
	method := reply.GetHttpMethod()
	handler := route.matcher.getActionHandler(path, method)
	if handler == nil {
		reply.SetStatusCode(404).With("404:you are lost")
		logger.Error("Not found Handler for[%s] [%s]", method, path)
		return
	}
	handler.doAction(reply)
}
