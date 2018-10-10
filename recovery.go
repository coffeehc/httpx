package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RecoveryMiddleware(logger *zap.Logger,devModule bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if logger != nil {
					logger.Error("HTTP请求异常", LogKeyHTTPContext(c)...)
				}
				if devModule{
					c.AbortWithStatusJSON(http.StatusInternalServerError,err)
				}else{
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		c.Next()
	}
}
