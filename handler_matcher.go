package web

import (
	"errors"
	"github.com/coffeehc/logger"
	"sort"
	"strings"
)

type handlerMatcher struct {
	requestHandlerMap map[HttpMethod]requestHandlerList
}

func (this *handlerMatcher) regeditAction(path string, method HttpMethod, requestHandler RequestHandler) error {
	newActionHandler, err := buildRequestHandler(path, method, requestHandler)
	if err != nil {
		logger.Error("添加Action出错:%s", err)
		return err
	}
	actionHandlers, ok := this.requestHandlerMap[method]
	if !ok {
		actionHandlers = make(requestHandlerList, 0)
	}
	for _, handler := range actionHandlers {
		if handler.exp.String() == newActionHandler.exp.String() {
			return errors.New(logger.Error("定义的uri[%s]与[%s]产生冲突,不能添加", handler.defineUri, newActionHandler.defineUri))
		}
	}
	this.requestHandlerMap[method] = append(actionHandlers, newActionHandler)
	logger.Debug("添加[%s] [%s] 对应的 Handler {%#v}", method, path, requestHandler)
	return nil
}

func (this *handlerMatcher) getActionHandler(uri string, method HttpMethod) *requestHandler {
	actionHandlers, ok := this.requestHandlerMap[method]
	if !ok {
		logger.Error("没有注册对应的Handler")
		return nil
	}
	paths := strings.Split(uri, PATH_SEPARATOR)
	pathSize := len(paths)
	for _, handler := range actionHandlers {
		if handler.pathSize == pathSize && handler.match(uri) {
			return handler
		}
	}
	return nil
}

func (this *handlerMatcher) sort() {
	for _, actionHandlers := range this.requestHandlerMap {
		sort.Sort(actionHandlers)
	}
}
