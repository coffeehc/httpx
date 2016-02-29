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
					if handler, ok := _router.errorHandlers[httpErr.Code]; ok {
						reply.GetResponseWriter().WriteHeader(httpErr.Code)
						handler(request, httpErr, reply)
						return
					}
					reply.SetStatusCode(httpErr.Code)
					reply.With(httpErr.Message)
				}
			}()
			chain(request, reply)
		},
		filterChainFunc: _router.handle,
	}
	return _router
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
