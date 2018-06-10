package httpx

import (
	"fmt"
	"time"

	"git.xiagaogao.com/coffee/boot/errors"
	"go.uber.org/zap"
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
	errorService errors.Service
	logger *zap.Logger
}

func newFilterWarp(matcher URIPatternMatcher, actionFilter Filter,errorService errors.Service, logger *zap.Logger) *filterWarp {
	return &filterWarp{matcher: matcher, requestFilter: actionFilter,errorService:errorService,logger:logger}
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

func NewAccessLogFilter(errorService errors.Service, logger *zap.Logger) Filter {
	return func(reply Reply, chain FilterChain) {
		t1 := time.Now()
		defer func() {
			printAccessLog(t1, reply,errorService,logger)
			if err := recover(); err != nil {
				panic(err)
			}
		}()
		chain(reply)
	}
}

func printAccessLog(startTime time.Time, reply Reply,errorService errors.Service, logger *zap.Logger) {
	request := reply.GetRequest()
	logger.Info(fmt.Sprintf("%s\t%s\t%s\t%d", time.Since(startTime), request.RemoteAddr, request.URL.RequestURI(), reply.GetStatusCode()))
}
