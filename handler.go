package web

import (
	"errors"
	"regexp"
	"strings"

	"github.com/coffeehc/logger"
)

type RequestHandler func(reply Reply)

type requestHandler struct {
	method           HttpMethod
	definePath       string
	hasPathFragments bool
	uriConversions   map[string]int
	pathSize         int
	requestHandler   RequestHandler
	exp              *regexp.Regexp
	//用于排序后提高 request 命中率
	//accessCount int64
}

func (this *requestHandler) doAction(reply Reply) {
	path := reply.GetPath()
	if this.hasPathFragments {
		paths := strings.Split(path, PATH_SEPARATOR)
		if this.pathSize != len(paths) {
			panic(errors.New(logger.Error("需要解析的uri[%s]不匹配定义的uri[%s]", path, this.definePath)))
		}
		for name, index := range this.uriConversions {
			reply.PutPathFragment(name, paths[index])
		}

	}
	this.requestHandler(reply)
}

//用于匹配是否
func (this *requestHandler) match(uri string) bool {
	if this.hasPathFragments {
		return this.exp.MatchString(uri)
	}
	return this.definePath == uri
}

func buildRequestHandler(path string, method HttpMethod, requestHandlerFunc RequestHandler) (*requestHandler, error) {
	if !strings.HasPrefix(path, PATH_SEPARATOR) {
		return nil, errors.New(logger.Error("定义的Uri必须是%s前缀", PATH_SEPARATOR))
	}
	paths := strings.Split(path, PATH_SEPARATOR)
	uriConversions := make(map[string]int, 0)
	conversionUri := ""
	pathSize := 0
	for index, p := range paths {
		pathSize++
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, WILDCARD_PREFIX) && strings.HasSuffix(p, WILDCARD_SUFFIX) {
			name := string([]byte(p)[len(WILDCARD_PREFIX) : len(p)-len(WILDCARD_SUFFIX)])
			uriConversions[name] = index
			p = _conversion
		}
		conversionUri += (PATH_SEPARATOR + p)
	}
	if conversionUri == "" {
		conversionUri = PATH_SEPARATOR
	}
	if strings.HasSuffix(path, PATH_SEPARATOR) {
		conversionUri += PATH_SEPARATOR
	}
	exp, err := regexp.Compile("^" + conversionUri + "$")
	if err != nil {
		return nil, err
	}
	return &requestHandler{
		exp:              exp,
		method:           method,
		definePath:       path,
		hasPathFragments: len(uriConversions) > 0,
		uriConversions:   uriConversions,
		pathSize:         pathSize,
		requestHandler:   requestHandlerFunc,
	}, nil
}

type requestHandlerList []*requestHandler

func (this requestHandlerList) Len() int {
	return len(this)
}
func (this requestHandlerList) Less(i, j int) bool {
	uri1s := strings.Split(this[i].exp.String(), PATH_SEPARATOR)
	uri2s := strings.Split(this[j].exp.String(), PATH_SEPARATOR)
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
func (this requestHandlerList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
