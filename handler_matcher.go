package httpx

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"git.xiagaogao.com/coffee/boot/errors"
	"git.xiagaogao.com/coffee/boot/logs"
	"go.uber.org/zap"
)

type handlerMatcher struct {
	requestHandlerMap map[RequestMethod]requestHandlerList
	errorService errors.Service
	logger *zap.Logger
}

func (impl *handlerMatcher) regeditAction(path string, method RequestMethod, requestHandler RequestHandler) errors.Error {
	newActionHandler, err := impl.buildRequestHandler(path, method, requestHandler)
	if err != nil {
		impl.logger.Error("添加Action出错", logs.F_Error(err))
		return err
	}
	actionHandlers, ok := impl.requestHandlerMap[method]
	if !ok {
		actionHandlers = make(requestHandlerList, 0)
	}
	for _, handler := range actionHandlers {
		if handler.exp.String() == newActionHandler.exp.String() {
			return impl.errorService.SystemError(fmt.Sprintf("定义的uri[%s]与[%s]产生冲突,不能添加", handler.defineURI, newActionHandler.defineURI))
		}
	}
	impl.requestHandlerMap[method] = append(actionHandlers, newActionHandler)
	impl.logger.Debug(fmt.Sprintf("添加[%s] [%s] 对应的 Handler {%#v}", method, path, requestHandler))
	return nil
}

func (impl *handlerMatcher) getActionHandler(uri string, method RequestMethod) *requestHandler {
	actionHandlers, ok := impl.requestHandlerMap[method]
	if !ok {
		impl.logger.Error("没有注册对应的Handler")
		return nil
	}
	paths := strings.Split(uri, pathSeparator)
	pathSize := len(paths)
	for _, handler := range actionHandlers {
		if handler.pathSize == pathSize && handler.match(uri) {
			return handler
		}
	}
	return nil
}

func (impl *handlerMatcher) sort() {
	for _, actionHandlers := range impl.requestHandlerMap {
		sort.Sort(actionHandlers)
	}
}


func (impl *handlerMatcher)buildRequestHandler(path string, method RequestMethod, requestHandlerFunc RequestHandler) (*requestHandler, errors.Error) {
	if !strings.HasPrefix(path, pathSeparator) {
		return nil, impl.errorService.SystemError(fmt.Sprintf("定义的Uri必须是%s前缀", pathSeparator))
	}
	u, err := url.ParseRequestURI(path)
	if err != nil {
		return nil, impl.errorService.SystemError("url格式错误",logs.F_ExtendData(path))
	}
	paths := strings.Split(u.Path, pathSeparator)
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
		return nil, impl.errorService.WrappedSystemError(err)
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