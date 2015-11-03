// define
package rest

type Error struct {
	Code             int32  `json:"code"`
	Debug_id         int64  `json:"debug_id"`
	Message          string `json:"message"`
	Information_link string `json:"information_link"`
}

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

func NewErrorResponse(errs ...Error) ErrorResponse {
	return ErrorResponse{errs}
}
