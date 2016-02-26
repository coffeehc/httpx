package web

import (
	"errors"
	"github.com/coffeehc/logger"
	"sort"
	"strings"
)

type actionMatcher struct {
	actionHandlerMap map[HttpMethod]actionHandlerList
}

func (this *actionMatcher) regeditAction(path string, method HttpMethod, requestHandler RequestHandler) error {
	newActionHandler, err := buildActionHandler(path, method, requestHandler)
	if err != nil {
		logger.Error("添加Action出错:%s", err)
		return err
	}
	actionHandlers, ok := this.actionHandlerMap[method]
	if !ok {
		actionHandlers = make(actionHandlerList, 0)
	}
	for _, handler := range actionHandlers {
		if handler.exp.String() == newActionHandler.exp.String() {
			return errors.New(logger.Error("定义的uri[%s]与[%s]产生冲突,不能添加", handler.defineUri, newActionHandler.defineUri))
		}
	}
	this.actionHandlerMap[method] = append(actionHandlers, newActionHandler)
	logger.Debug("添加[%s] [%s] 对应的 Handler {%#t}", method, path, requestHandler)
	return nil
}

func (this *actionMatcher) getActionHandler(uri string, method HttpMethod) *actionHandler {
	actionHandlers, ok := this.actionHandlerMap[method]
	if !ok {
		logger.Error("没有注册对应的Handler")
		return nil
	}
	//logger.Debug("%#q",actionHandlers)
	paths := strings.Split(uri, PATH_SEPARATOR)
	pathSize := len(paths)
	for _, handler := range actionHandlers {
		if handler.pathSize == pathSize && handler.match(uri) {
			return handler
		}
	}
	return nil
}

func (this *actionMatcher) sort() {
	for _, actionHandlers := range this.actionHandlerMap {
		sort.Sort(actionHandlers)
	}
}
