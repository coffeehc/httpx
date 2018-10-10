package httpx

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AccessLogMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		logger.Debug("HTTP访问", LogKeyHTTPContext(c, zap.Duration("http_delay", time.Since(t)))...)
	}
}
