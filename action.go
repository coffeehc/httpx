package web

import (
	"errors"
	"github.com/coffeehc/logger"
	"net/http"
	"regexp"
	"strings"
)

type RequestHandler func(request *http.Request, pathFragments map[string]string, reply *Reply)

type actionHandler struct {
	method           HttpMethod
	defineUri        string
	hasPathFragments bool
	uriConversions   map[string]int
	pathSize         int
	requestHandle    func(request *http.Request, PathFragments map[string]string, reply *Reply)
	exp              *regexp.Regexp
	//用于排序后提高 action 命中率
	//accessCount int64
}

func (this *actionHandler) doAction(request *http.Request, reply *Reply) {
	requestUri := request.URL.Path
	param := make(map[string]string, 0)
	if this.hasPathFragments {
		paths := strings.Split(requestUri, PATH_SEPARATOR)
		if this.pathSize != len(paths) {
			panic(errors.New(logger.Error("需要解析的uri[%s]不匹配定义的uri[%s]", requestUri, this.defineUri)))
		}
		for name, index := range this.uriConversions {
			param[name] = paths[index]
		}
	}
	this.requestHandle(request, param, reply)
}

//用于匹配是否
func (this *actionHandler) match(uri string) bool {
	return this.exp.MatchString(uri)
}

func buildActionHandler(path string, method HttpMethod, requestHandler RequestHandler) (*actionHandler, error) {
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
	exp, err := regexp.Compile("^" + conversionUri + "$")
	if err != nil {
		return nil, err
	}
	return &actionHandler{
		exp:              exp,
		method:           method,
		defineUri:        path,
		hasPathFragments: len(uriConversions) > 0,
		uriConversions:   uriConversions,
		pathSize:         pathSize,
		requestHandle:    requestHandler,
	}, nil
}

type actionHandlerList []*actionHandler

func (this actionHandlerList) Len() int {
	return len(this)
}
func (this actionHandlerList) Less(i, j int) bool {
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
func (this actionHandlerList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
