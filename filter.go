package httpx

import (
	"time"

	"github.com/coffeehc/logger"
)

//Filter 定义实现 Filter 需要实现的方法类型
type Filter func(reply Reply, chain FilterChain)

//FilterChain 定义的需要执行的方法子方法
type FilterChain func(reply Reply)

type filterWarp struct {
	matcher         URIPatternMatcher
	requestFilter   Filter
	nextFilter      *filterWarp
	filterChainFunc FilterChain
}

func newFilterWarp(matcher URIPatternMatcher, actionFilter Filter) *filterWarp {
	return &filterWarp{matcher: matcher, requestFilter: actionFilter}
}

func (warp *filterWarp) addNextFilter(filter *filterWarp) {
	if warp.nextFilter == nil {
		filter.filterChainFunc = warp.filterChainFunc
		warp.nextFilter = filter
		warp.filterChainFunc = nil
		return
	}
	warp.nextFilter.addNextFilter(filter)
}

func (warp *filterWarp) doFilter(reply Reply) {
	path := reply.GetRequest().URL.Path
	if warp.matcher.match(path) {
		warp.requestFilter(reply, warp.filterChain)
		return
	}
	warp.filterChain(reply)
}

func (warp *filterWarp) filterChain(reply Reply) {
	if warp.nextFilter == nil {
		warp.filterChainFunc(reply)
		return
	}
	warp.nextFilter.doFilter(reply)
}

//AccessLogFilter 访问日志的 Filter, 在 consol 上输出每个被访问的 url 及响应时间
func AccessLogFilter(reply Reply, chain FilterChain) {
	t1 := time.Now()
	defer func() {
		printAccessLog(t1, reply)
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	chain(reply)
}

func printAccessLog(startTime time.Time, reply Reply) {
	request := reply.GetRequest()
	logger.Info("%s\t%s\t%s\t%d", time.Since(startTime), request.RemoteAddr, request.URL.RequestURI(), reply.GetStatusCode())
}
