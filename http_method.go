package web

import "net/http"

//RequestMethod 请求的 Http 方法
type RequestMethod string

const (
	//GET GET Method
	GET = RequestMethod(http.MethodGet)
	//HEAD HEAD Method
	HEAD = RequestMethod(http.MethodHead)
	//POST POST Method
	POST = RequestMethod(http.MethodPost)
	//PUT PUT Method
	PUT = RequestMethod(http.MethodPut)
	//PATCH PATCH Method
	PATCH = RequestMethod(http.MethodPatch)
	//DELETE DELETE Method
	DELETE = RequestMethod(http.MethodDelete)
	//CONNECT CONNECT Method
	CONNECT = RequestMethod(http.MethodConnect)
	//OPTIONS OPTIONS Method
	OPTIONS = RequestMethod(http.MethodOptions)
	//TRACE TRACE Method
	TRACE = RequestMethod(http.MethodTrace)
)
