package httpx

import (
	"regexp"
	"strings"

	"net/url"

	"git.xiagaogao.com/coffee/boot/errors"
	"git.xiagaogao.com/coffee/boot/logs"
	"go.uber.org/zap"
)

//RequestHandler 处理 Request的方法定义
type RequestHandler func(reply Reply)

type requestHandler struct {
	method            RequestMethod
	defineURI         string
	hasPathFragments  bool
	uriConversions    map[string]int
	pathSize          int
	requestHandleFunc RequestHandler
	exp               *regexp.Regexp
	//用于排序后提高 request 命中率
	//accessCount int64
	errorService errors.Service
	logger *zap.Logger
}

func (handler *requestHandler) doAction(reply Reply) {
	requestURI := reply.GetRequest().RequestURI
	if handler.hasPathFragments {
		u, err := url.ParseRequestURI(requestURI)
		if err != nil {
			handler.logger.Error("错误的uri", logs.F_ExtendData(requestURI))
			reply.SetStatusCode(500)
			return
		}
		paths := strings.Split(u.Path, pathSeparator)
		if handler.pathSize != len(paths) {
			reply.SetStatusCode(404)
			return
		}
		for name, index := range handler.uriConversions {
			reply.AddPathFragment(name, paths[index])
		}
	}
	handler.requestHandleFunc(reply)
}

//用于匹配是否
func (handler *requestHandler) match(uri string) bool {
	if handler.hasPathFragments {
		return handler.exp.MatchString(uri)
	}
	return handler.defineURI == uri
}



type requestHandlerList []*requestHandler

func (list requestHandlerList) Len() int {
	return len(list)
}
func (list requestHandlerList) Less(i, j int) bool {
	h1 := list[i]
	h2 := list[j]
	if !h1.hasPathFragments && h2.hasPathFragments {
		return true
	}
	if !h2.hasPathFragments && h1.hasPathFragments {
		return false
	}
	uri1s := strings.Split(h1.exp.String(), pathSeparator)
	uri2s := strings.Split(h2.exp.String(), pathSeparator)
	if len(uri1s) != len(uri2s) {
		return len(uri1s) > len(uri2s)
	}
	for i, path1 := range uri1s {
		path2 := uri2s[i]
		if path1 != path2 {
			if path2 == _conversion {
				return false
			}
			return len(path1) < len(path2)
		}
	}
	return true
}
func (list requestHandlerList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}
