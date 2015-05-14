// filter
package web

import (
	"net/http"
	"time"

	"github.com/coffeehc/logger"
)

type ActionFilter func(request *http.Request, reply *Reply, chain FilterChain)

type FilterChain func(request *http.Request, reply *Reply)

type filterWarp struct {
	matcher      uriPatternMatcher
	index        int
	actionFilter ActionFilter
	dispatcher   *routingDispatcher
}

func (this *filterWarp) filter(request *http.Request, reply *Reply) {
	this.actionFilter(request, reply, this.filterChain)
}
func (this *filterWarp) filterChain(request *http.Request, reply *Reply) {
	chain := this.dispatcher.nextFilter(request, this.index)
	if chain == nil {
		this.dispatcher.handle(request, reply)
		return
	}
	chain.filter(request, reply)
}
func newFilterWarp(matcher uriPatternMatcher, index int, actionFilter ActionFilter, dispatcher *routingDispatcher) *filterWarp {
	return &filterWarp{matcher: matcher, index: index, actionFilter: actionFilter, dispatcher: dispatcher}
}

func AccessLogFilter(request *http.Request, reply *Reply, chain FilterChain) {
	t1 := time.Now()
	chain(request, reply)
	delay := time.Since(t1)
	logger.Info("访问时间[%s]处理时长[%s]地址[%s]返回状态[%d]", t1, delay, request.RemoteAddr, reply.GetStatusCode())
}
