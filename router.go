package httpx

import (
	"fmt"

	"strings"

	"git.xiagaogao.com/coffee/boot/errors"
	"git.xiagaogao.com/coffee/boot/logs"
	"go.uber.org/zap"
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
	errorService errors.Service
	logger *zap.Logger
}

func newRouter(errorService errors.Service,logger *zap.Logger) *router {
	_router := &router{matcher: handlerMatcher{requestHandlerMap: make(map[RequestMethod]requestHandlerList),errorService:errorService,logger:logger}}
	_router.errorHandlers = errorHandlers(make(map[int]RequestErrorHandler, 0))
	_router.filter = &filterWarp{
		matcher: newServletStyleURIPatternMatcher("/*",logger),
		requestFilter: func(reply Reply, chain FilterChain) {
			defer func() {
				if err := recover(); err != nil {
					var httpErr *HTTPError
					var ok bool
					if httpErr, ok = err.(*HTTPError); !ok {
						httpErr = NewHTTPErr(500, fmt.Sprintf("%s", err))
					}
					logger.Error("http err", logs.F_ExtendData(err))
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
	newFilter := newFilterWarp(matcher, actionFilter,r.errorService,r.logger)
	newFilter.nextFilter = oldFilter
	r.filter = newFilter

}

func (r *router) addLastFilter(matcher URIPatternMatcher, actionFilter Filter) {
	r.filter.addNextFilter(newFilterWarp(matcher, actionFilter,r.errorService,r.logger))
}

func (r *router) handle(reply Reply){
	request := reply.GetRequest()
	request.ParseForm()
	request.URL.Path = strings.Replace(request.URL.Path, "//", "/", -1)
	method := RequestMethod(strings.ToUpper(request.Method))
	handler := r.matcher.getActionHandler(request.URL.Path, method)
	if handler == nil {
		reply.SetStatusCode(404).With(NewHTTPErr(404, "you are lost")).As(DefaultRenderJSON)
		r.logger.Error(fmt.Sprintf("Not found Handler for[%s] [%s]", method, request.URL.Path))
		return
	}
	handler.doAction(reply)
}
