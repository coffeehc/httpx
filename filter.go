// filter
package web

import (
	"net/http"
	"time"

	"github.com/coffeehc/logger"
)

type Filter func(request *http.Request, reply Reply, chain FilterChain)

type FilterChain func(request *http.Request, reply Reply)

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

func (this *filterWarp) doFilter(request *http.Request, reply Reply) {
	if this.matcher.match(request.URL.Path) {
		this.requestFilter(request, reply, this.filterChain)
		return
	}
	this.filterChain(request, reply)
}

func (this *filterWarp) filterChain(request *http.Request, reply Reply) {
	if this.nextFilter == nil {
		this.filterChain(request, reply)
		return
	}
	this.nextFilter.doFilter(request, reply)
}

func AccessLogFilter(request *http.Request, reply Reply, chain FilterChain) {
	t1 := time.Now()
	chain(request, reply)
	delay := time.Since(t1)
	logger.Info("%s\t%s\t%s\t%s\t%d", t1, delay, request.RemoteAddr, request.URL, reply.GetStatusCode())
}
