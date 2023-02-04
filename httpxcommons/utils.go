package httpxcommons

import (
	"github.com/coffeehc/base/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"
)

func SendSuccess(c *fiber.Ctx, obj interface{}, code int64) error {
	if c.Accepts("application/x-protobuf") != "" {
		//log.Debug("+++", zap.String("Accepts", c.Accepts("application/x-protobuf")))
		var data []byte
		var err error
		message, ok := obj.(string)
		if !ok {
			msg, ok := obj.(proto.Message)
			if !ok {
				//log.Error("========")
				c.SendStatus(500)
				return nil
			}
			data, err = proto.Marshal(msg)
			if err != nil {
				//log.Error("========>>>>")
				c.SendStatus(500)
				return nil
			}
		} else {
			data = []byte(message)
		}
		resp := &PBResponse{
			Code:    code,
			Success: true,
			Payload: data,
		}
		data, err = proto.Marshal(resp)
		if err != nil {
			log.Error("错误", zap.Error(err))
			c.SendStatus(500)
			return nil
		}
		c.Status(200).BodyParser(data)
		return nil
	}
	return c.JSON(&AjaxResponse{
		Code:    code,
		Payload: obj,
		Success: true,
	})
}

func SendErrorWithRedirect(c *fiber.Ctx, message string, redirect string, code int64, statusCode int) error {
	if c.Accepts("application/x-protobuf") != "" {
		resp := &PBResponse{
			Code:    code,
			Message: message,
		}
		data, err := proto.Marshal(resp)
		if err != nil {
			log.Error("错误", zap.Error(err))
			return c.SendStatus(500)
		}
		return c.Status(statusCode).Send(data)

	}
	return c.Status(statusCode).JSON(&AjaxResponse{
		Code:     code,
		Message:  message,
		Redirect: redirect,
	})
} //(c, "", "/user/login", 401, 401)

func SendError(c *fiber.Ctx, err string, code int64, statusCode int) error {
	if c.Accepts("application/x-protobuf") != "" {
		resp := &PBResponse{
			Code:    code,
			Message: err,
		}
		data, err := proto.Marshal(resp)
		if err != nil {
			log.Error("错误", zap.Error(err))
			return c.SendStatus(500)

		}
		return c.Status(statusCode).Send(data)

	}
	return c.Status(statusCode).JSON(&AjaxResponse{
		Code:    code,
		Message: err,
	})
}

func SendErrors(c *fiber.Ctx, err error, code int64, statusCode int) error {
	if c.Accepts("application/x-protobuf") != "" {
		resp := &PBResponse{
			Code:    code,
			Message: err.Error(),
		}
		data, err := proto.Marshal(resp)
		if err != nil {
			log.Error("错误", zap.Error(err))
			return c.SendStatus(500)
		}
		return c.Status(statusCode).Send(data)
	}
	return c.Status(statusCode).JSON(&AjaxResponse{
		Code:    code,
		Message: err.Error(),
	})
}
