package web

import "net/http"

type HttpMethod string

const (
	GET     = HttpMethod(http.MethodGet)
	HEAD    = HttpMethod(http.MethodHead)
	POST    = HttpMethod(http.MethodPost)
	PUT     = HttpMethod(http.MethodPut)
	PATCH   = HttpMethod(http.MethodPatch)
	DELETE  = HttpMethod(http.MethodDelete)
	CONNECT = HttpMethod(http.MethodConnect)
	OPTIONS = HttpMethod(http.MethodOptions)
	TRACE   = HttpMethod(http.MethodTrace)
)
