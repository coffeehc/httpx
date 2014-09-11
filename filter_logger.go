// filter_recover
package web

import (
	"net/http"
	"time"

	"github.com/coffeehc/logger"
)

type LoggerFilter struct {
	Pattern string
}

func (this *LoggerFilter) GetPattern() string {
	return this.Pattern
}

func (this *LoggerFilter) Init(conf *WebConfig) {
}

func (this *LoggerFilter) DoFilter(req *http.Request, reply *Reply, chain FilterChain) {
	path := req.URL.Path
	startTime := time.Now()
	chain.DoFilter(req, reply)
	endTime := time.Now()
	logger.Debug("处理请求[%s]耗时:%4.3f毫秒", path, float64(endTime.Sub(startTime))/float64(time.Millisecond))
}

func (this *LoggerFilter) Destroy() {
}
