package httpx

import (
	"context"
	es "errors"
	"github.com/coffeehc/base/errors"
	"github.com/coffeehc/base/log"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"strings"
	"time"
)

func AccessLogMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		t := time.Now()
		e := c.Next()
		if e == nil {
			log.Debug(c.Path(), zap.Duration("times", time.Since(t)), zap.String("method", c.Method()))
		} else {
			log.Error("HTTP访问", zap.Duration("times", time.Since(t)), zap.String("mothod", c.Method()), zap.String("path", c.Path()), zap.Int("status", c.Response().StatusCode()))
		}
		return e
	}
}

func RecoverMiddleware(t time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		timeoutContext, cancel := context.WithTimeout(c.UserContext(), t)
		defer func() {
			cancel()
			if err := recover(); err != nil {
				path := c.Path()
				if errStr, ok := err.(string); ok {
					log.DPanic(errStr, zap.String("method", c.Method()), zap.String("path", path))
				} else {
					e := errors.ConverUnknowError(err)
					if errors.IsMessageError(e) {
						log.Error(e.Error(), zap.String("method", c.Method()), zap.String("path", path))
					} else {
						if strings.HasPrefix(e.Error(), "context ") || strings.HasPrefix(e.Error(), "rpc error") {
							log.Error(e.Error(), zap.String("method", c.Method()), zap.String("path", path))
						} else {
							log.DPanic(e.Error(), zap.String("method", c.Method()), zap.String("path", path))
						}
					}
				}
				c.SendStatus(500)
			}
		}()
		c.SetUserContext(timeoutContext)
		err := c.Next()
		if err != nil {
			if es.Is(err, context.DeadlineExceeded) {
				return fiber.ErrRequestTimeout
			}
		}
		return err
	}
}
