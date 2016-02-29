package websocket

import (
	"github.com/coffeehc/logger"
	"github.com/coffeehc/web"
)

func (server *web.Server) RegeditWebSocket(path string, service WebScoketHandler) error {
	adapter := &webScoketAdapter{service}
	err := server.router.matcher.regeditAction(path, web.GET, adapter.webScoketHandlerAdapter)
	if err != nil {
		logger.Error("注册 Handler 失败:%s", err)
	}
	return err
}
