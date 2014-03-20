// filter.go
package web

import "net/http"
import "regexp"
import "fmt"

type Filter interface {
	DoFilter(req *http.Request, reply *Reply, chain FilterChain) error
}

type FilterChain interface {
	DoFilter(req *http.Request, reply *Reply) error
}

type filterChainInvocation struct {
	index int //-1
	conf  *WebConfig
}

func newFilterChainInvocation(conf *WebConfig) *filterChainInvocation {
	fci := new(filterChainInvocation)
	fci.conf = conf
	fci.index = -1
	return fci
}

func (this *filterChainInvocation) DoFilter(req *http.Request, reply *Reply) (err error) {
	this.index++
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			//TODO 捕获异常后的处理
			err = fmt.Errorf("处理异常:%s", recoverErr)
			return
		}
	}()
	if this.index < len(filterDefinitions) {
		err = filterDefinitions[this.index].doFilter(req, reply, this)
	} else {
		request := &Request{Request: req}
		dispatcher.dispatch(request, reply)
		//TODO 这里其实可以插入一个缓冲器,保证在一定时间内访问过高的uri进行缓冲
	}
	return
}

type filterDefinition struct {
	pattern *regexp.Regexp
	filter  Filter
}

func (this *filterDefinition) doFilter(req *http.Request, reply *Reply, chainInvocation *filterChainInvocation) (err error) {
	if this.pattern.MatchString(req.URL.Path) {
		err = this.filter.DoFilter(req, reply, chainInvocation)
	} else {
		err = chainInvocation.DoFilter(req, reply)
	}
	return
}

var filterDefinitions []*filterDefinition

func init() {
	filterDefinitions = make([]*filterDefinition, 0)
}

func AddFilter(pattern string, filter Filter) {
	fd := new(filterDefinition)
	r, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("Filter指定的UriPattern不能解析:%s", pattern))
	}
	fd.pattern = r
	fd.filter = filter
	filterDefinitions = append(filterDefinitions, fd)
}
