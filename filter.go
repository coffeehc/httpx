// filter
package web

import (
	"net/http"
	"time"

	"github.com/coffeehc/logger"
	"net"
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
		this.filterChainFunc(request, reply)
		return
	}
	this.nextFilter.doFilter(request, reply)
}

func SimpleAccessLogFilter(request *http.Request, reply Reply, chain FilterChain) {
	t1 := time.Now()
	defer func() {
		if err := recover(); err != nil {
			if httpError, ok := err.(*HttpError); ok {
				reply.SetStatusCode(httpError.Code)
			}
			pringAccessLog(t1, request, reply)
			panic(err)
		}
	}()
	chain(request, reply)
	pringAccessLog(t1, request, reply)
}

func pringAccessLog(startTime time.Time, request *http.Request, reply Reply) {
	delay := time.Since(startTime)
	addr, err := net.ResolveTCPAddr("tcp", request.RemoteAddr)
	var ip string
	if err != nil {
		ip = "0.0.0.0"
	} else {
		ip = addr.IP.String()
	}
	logger.Info("%s\t%s\t%s\t%d", delay, ip, request.URL, reply.GetStatusCode())
}
