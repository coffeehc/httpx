package web

import (
	"fmt"
)

// HTTPError 对 http 处理的过程中异常的封装
type HTTPError struct {
	Code    int         `json:"httpcode"`
	Message interface{} `json:"message"`
}

//Error 实现的 error接口
func (httpError *HTTPError) Error() string {
	return fmt.Sprintf("%d:%s", httpError.Code, httpError.Message)
}

//NewHTTPErr 创建一个新的 Http 错误描述
func NewHTTPErr(code int, message string) *HTTPError {
	return &HTTPError{code, message}
}

// TODO 需要有一个地方能够注册异常的处理
type errorHandlers map[int]RequestErrorHandler

//RequestErrorHandler 请求异常的处理方法定义
type RequestErrorHandler func(err *HTTPError, reply Reply)
