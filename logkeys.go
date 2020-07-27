package httpx

import (
  "github.com/gin-gonic/gin"
  "go.uber.org/zap"
)

func LogKeyHTTPContext(c *gin.Context, field ...zap.Field) []zap.Field {
  baseFileds := []zap.Field{zap.Int("statusCode", c.Writer.Status()), zap.String("http_client_ip", c.ClientIP()), zap.String("http_method", c.Request.Method), zap.String("http_path", c.Request.URL.Path)}
  if len(field) > 0 {
    return append(field, baseFileds...)
  }
  return baseFileds
}
