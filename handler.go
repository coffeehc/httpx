package httpx

import (
	"errors"
	"regexp"
	"strings"

	"github.com/coffeehc/logger"
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
}

func (handler *requestHandler) doAction(reply Reply) {
	requestURI := reply.GetRequest().RequestURI
	if handler.hasPathFragments {
		paths := strings.Split(requestURI, pathSeparator)
		if handler.pathSize != len(paths) {
			panic(errors.New(logger.Error("需要解析的uri[%s]不匹配定义的uri[%s]", requestURI, handler.defineURI)))
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

func buildRequestHandler(path string, method RequestMethod, requestHandlerFunc RequestHandler) (*requestHandler, error) {
	if !strings.HasPrefix(path, pathSeparator) {
		return nil, errors.New(logger.Error("定义的Uri必须是%s前缀", pathSeparator))
	}
	paths := strings.Split(path, pathSeparator)
	uriConversions := make(map[string]int, 0)
	conversionURI := ""
	pathSize := 0
	for index, p := range paths {
		pathSize++
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, wildcardPrefix) && strings.HasSuffix(p, wildcardSuffix) {
			name := string([]byte(p)[len(wildcardPrefix) : len(p)-len(wildcardSuffix)])
			uriConversions[name] = index
			p = _conversion
		}
		conversionURI += (pathSeparator + p)
	}
	if conversionURI == "" {
		conversionURI = pathSeparator
	}
	if strings.HasSuffix(path, pathSeparator) {
		conversionURI += pathSeparator
	}
	exp, err := regexp.Compile("^" + conversionURI + "$")
	if err != nil {
		return nil, err
	}
	return &requestHandler{
		exp:               exp,
		method:            method,
		defineURI:         path,
		hasPathFragments:  len(uriConversions) > 0,
		uriConversions:    uriConversions,
		pathSize:          pathSize,
		requestHandleFunc: requestHandlerFunc,
	}, nil
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
