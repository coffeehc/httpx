package web

import (
	"fmt"
)

type HttpError struct {
	Code    int         `json:"httpcode"`
	Message interface{} `json:"message"`
}

func (this *HttpError) Error() string {
	return fmt.Sprintf("%d:%s", this.Code, this.Message)
}

func NewHttpErr(code int, message string) *HttpError {
	return &HttpError{code, message}
}

func HTTPERR_500(message string) *HttpError {
	return &HttpError{500, message}
}

func HTTPERR_400(message string) *HttpError {
	return &HttpError{400, message}
}

type ErrorHandlers map[int]RequestErrorHandler

type RequestErrorHandler func(err *HttpError, reply Reply)
