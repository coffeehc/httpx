package httpxcommons

import (
	"github.com/coffeehc/base/errors"
	"github.com/coffeehc/base/log"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"strings"
)

func SendPBSuccess(c *fiber.Ctx, obj interface{}, code int64) error {
	msg, ok := obj.(proto.Message)
	if !ok {
		//log.Error("========")
		c.SendStatus(fiber.StatusNotAcceptable)
		return nil
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		//log.Error("========>>>>")
		return err
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
	return c.Status(200).Send(data)
}

func SendSuccess(c fiber.Ctx, obj interface{}, code int64) error {
	if !strings.Contains(c.Get(fiber.HeaderAccept), "*/*") && c.Accepts("application/x-protobuf") != "" {
		//log.Debug("+++", zap.String("Accepts", c.Accepts("application/x-protobuf")))
		var data []byte
		var err error
		message, ok := obj.(string)
		if !ok {
			msg, ok := obj.(proto.Message)
			if !ok {
				//log.Error("========")
				c.SendStatus(fiber.StatusNotAcceptable)
				return nil
			}
			data, err = proto.Marshal(msg)
			if err != nil {
				//log.Error("========>>>>")
				c.SendStatus(501)
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
			return c.SendStatus(500)
		}
		return c.Status(200).Send(data)
	}
	return c.JSON(&AjaxResponse{
		Code:    code,
		Payload: obj,
		Success: true,
	})
}

func SendErrorWithRedirect(c fiber.Ctx, message string, redirect string, code int64, statusCode int) error {
	if !strings.Contains(c.Get(fiber.HeaderAccept), "*/*") && c.Accepts("application/x-protobuf") != "" {
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

func SendError(c fiber.Ctx, err string, code int64, statusCode int) error {
	if !strings.Contains(c.Get(fiber.HeaderAccept), "*/*") && c.Accepts("application/x-protobuf") != "" {
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

func SendErrors(c fiber.Ctx, err error, code int64, statusCode int) error {
	message := err.Error()
	if errors.IsSystemError(err) || errors.IsDBError(err) {
		message = "系统内部错误"
	}
	if !strings.Contains(c.Get(fiber.HeaderAccept), "*/*") && c.Accepts("application/x-protobuf") != "" {
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
		Code:    code,
		Message: message,
	})
}
