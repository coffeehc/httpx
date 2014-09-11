// filter_recover
package web

import (
	"net/http"

	"github.com/coffeehc/logger"
)

type RecoverFilter struct {
	Pattern string
}

func (this *RecoverFilter) GetPattern() string {
	return this.Pattern
}

func (this *RecoverFilter) Init(conf *WebConfig) {
}

func (this *RecoverFilter) DoFilter(req *http.Request, reply *Reply, chain FilterChain) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			logger.Error("请求[%s]出现异常:%s", req.URL, recoverErr)
			reply.Error("出现了不可恢复的异常", 500)
		}
	}()
	chain.DoFilter(req, reply)
}

func (this *RecoverFilter) Destroy() {
}
