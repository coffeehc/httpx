// filter.go
package web

import "net/http"
import "regexp"
import "fmt"
import "github.com/coffeehc/logger"

type Filter interface {
	DoFilter(req *http.Request, w http.ResponseWriter, chain FilterChain) error
}

type FilterChain interface {
	DoFilter(req *http.Request, w http.ResponseWriter) error
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

func (this *filterChainInvocation) DoFilter(req *http.Request, w http.ResponseWriter) error {
	this.index++
	if this.index < len(filterDefinitions) {
		err := filterDefinitions[this.index].doFilter(req, w, this)
		if err != nil {
			return err
		}
	} else {
		request := new(Request)
		request.Request = req
		//TODO 这里其实可以插入一个缓冲器,保证在一定时间内访问过高的uri进行缓冲
		reply := dispatcher.dispatch(request)
		if reply == nil {
			reply = NewReply(request).NoFindPage()
		}
		err := reply.writeResponse(w, req)
		if err != nil {
			logger.Errorf("出现一个错误%s", err)
			return err
		}
	}
	return nil
}

type filterDefinition struct {
	pattern *regexp.Regexp
	filter  Filter
}

func (this *filterDefinition) doFilter(req *http.Request, w http.ResponseWriter, chainInvocation *filterChainInvocation) error {
	var err error
	if this.pattern.MatchString(req.URL.Path) {
		err = this.filter.DoFilter(req, w, chainInvocation)
	} else {
		err = chainInvocation.DoFilter(req, w)
	}
	if err != nil {
		return err
	}
	return nil
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
