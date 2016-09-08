// filter
package web

import (
	"time"

	"github.com/coffeehc/logger"
)

type Filter func(reply Reply, chain FilterChain)

type FilterChain func(reply Reply)

type filterWarp struct {
	matcher         UriPatternMatcher
	requestFilter   Filter
	nextFilter      *filterWarp
	filterChainFunc FilterChain
}

func newFilterWarp(matcher UriPatternMatcher, actionFilter Filter) *filterWarp {
	return &filterWarp{matcher: matcher, requestFilter: actionFilter}
}

func (this *filterWarp) addNextFilter(filter *filterWarp) {
	if this.nextFilter == nil {
		filter.filterChainFunc = this.filterChainFunc
		this.nextFilter = filter
		this.filterChainFunc = nil
		return
	}
	this.nextFilter.addNextFilter(filter)
}

func (this *filterWarp) doFilter(reply Reply) {
	if this.matcher.match(reply.GetPath()) {
		this.requestFilter(reply, this.filterChain)
		return
	}
	this.filterChain(reply)
}

func (this *filterWarp) filterChain(reply Reply) {
	if this.nextFilter == nil {
		this.filterChainFunc(reply)
		return
	}
	this.nextFilter.doFilter(reply)
}

func SimpleAccessLogFilter(reply Reply, chain FilterChain) {
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
	logger.Info("%s\t%s\t%s\t%d", time.Since(startTime), reply.GetRemoteAddr().String(), reply.GetFullURL(), reply.GetStatusCode())
}
