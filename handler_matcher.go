package httpx

import (
	"errors"
	"sort"
	"strings"

	"github.com/coffeehc/logger"
)

type handlerMatcher struct {
	requestHandlerMap map[RequestMethod]requestHandlerList
}

func (hm *handlerMatcher) regeditAction(path string, method RequestMethod, requestHandler RequestHandler) error {
	newActionHandler, err := buildRequestHandler(path, method, requestHandler)
	if err != nil {
		logger.Error("添加Action出错:%s", err)
		return err
	}
	actionHandlers, ok := hm.requestHandlerMap[method]
	if !ok {
		actionHandlers = make(requestHandlerList, 0)
	}
	for _, handler := range actionHandlers {
		if handler.exp.String() == newActionHandler.exp.String() {
			return errors.New(logger.Error("定义的uri[%s]与[%s]产生冲突,不能添加", handler.defineURI, newActionHandler.defineURI))
		}
	}
	hm.requestHandlerMap[method] = append(actionHandlers, newActionHandler)
	logger.Debug("添加[%s] [%s] 对应的 Handler {%#v}", method, path, requestHandler)
	return nil
}

func (hm *handlerMatcher) getActionHandler(uri string, method RequestMethod) *requestHandler {
	actionHandlers, ok := hm.requestHandlerMap[method]
	if !ok {
		logger.Error("没有注册对应的Handler")
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

func (hm *handlerMatcher) sort() {
	for _, actionHandlers := range hm.requestHandlerMap {
		sort.Sort(actionHandlers)
	}
}
