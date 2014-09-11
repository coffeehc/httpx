// filter.go
package web

import (
	"inject"
	"net/http"
)
import "regexp"

type Filter interface {
	GetPattern() string
	Init(conf *WebConfig) //初始化
	DoFilter(req *http.Request, reply *Reply, chain FilterChain)
	Destroy() //销毁
}

type FilterChain interface {
	DoFilter(req *http.Request, reply *Reply)
}

type filterChainInvocation struct {
	filterDefinitions []*filterDefinition
	inject.Injector
	index int //-1
}

func newFilterChainInvocation(conf *WebConfig) *filterChainInvocation {
	fci := new(filterChainInvocation)
	fci.filterDefinitions = conf.filterDefinitions
	fci.Injector = inject.NewChild(conf.Injector)
	fci.index = -1
	return fci
}

func (this *filterChainInvocation) DoFilter(req *http.Request, reply *Reply) {
	this.index++
	if this.index < len(this.filterDefinitions) {
		this.filterDefinitions[this.index].doFilter(req, reply, this)
	} else {
		dispatcher.dispatch(req, reply)
		//TODO 这里其实可以插入一个缓冲器,保证在一定时间内访问过高的uri进行缓冲
	}
}

type filterDefinition struct {
	pattern *regexp.Regexp
	filter  Filter
}

func (this *filterDefinition) doFilter(req *http.Request, reply *Reply, chainInvocation *filterChainInvocation) {
	if this.pattern.MatchString(req.URL.Path) {
		this.filter.DoFilter(req, reply, chainInvocation)
	} else {
		chainInvocation.DoFilter(req, reply)
	}
	return
}
