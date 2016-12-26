package web

import (
	"fmt"

	"github.com/coffeehc/logger"
	"strings"
)

const (
	pathSeparator  = "/"
	wildcardPrefix = "{"
	wildcardSuffix = "}"
	_conversion    = "[^/]*"
)

type router struct {
	matcher       handlerMatcher
	filter        *filterWarp
	errorHandlers errorHandlers
}

func newRouter() *router {
	_router := &router{matcher: handlerMatcher{requestHandlerMap: make(map[RequestMethod]requestHandlerList)}}
	_router.errorHandlers = errorHandlers(make(map[int]RequestErrorHandler, 0))
	_router.filter = &filterWarp{
		matcher: newServletStyleURIPatternMatcher("/*"),
		requestFilter: func(reply Reply, chain FilterChain) {
			defer func() {
				if err := recover(); err != nil {
					var httpErr *HTTPError
					var ok bool
					if httpErr, ok = err.(*HTTPError); !ok {
						httpErr = NewHTTPErr(500, fmt.Sprintf("%#s", err))
					}
					reply.SetStatusCode(httpErr.Code)
					if handler, ok := _router.errorHandlers[httpErr.Code]; ok {
						handler(httpErr, reply)
						return
					}
					reply.With(httpErr.Message).As(DefaultRenderJSON)
				}
			}()
			chain(reply)
		},
		filterChainFunc: _router.handle,
	}
	return _router
}

func (r *router) addFirstFilter(matcher URIPatternMatcher, actionFilter Filter) {
	oldFilter := r.filter
	newFilter := newFilterWarp(matcher, actionFilter)
	newFilter.nextFilter = oldFilter
	r.filter = newFilter

}

func (r *router) addLastFilter(matcher URIPatternMatcher, actionFilter Filter) {
	r.filter.addNextFilter(newFilterWarp(matcher, actionFilter))
}

func (r *router) handle(reply Reply) {
	request := reply.GetRequest()
	request.ParseForm()
	request.URL.Path = strings.Replace(request.URL.Path, "//", "/", -1)
	path := request.RequestURI
	method := RequestMethod(strings.ToUpper(request.Method))
	handler := r.matcher.getActionHandler(path, method)
	if handler == nil {
		reply.SetStatusCode(404).With(NewHTTPErr(404, "you are lost")).As(DefaultRenderJSON)
		logger.Error("Not found Handler for[%s] [%s]", method, path)
		return
	}
	handler.doAction(reply)
}
