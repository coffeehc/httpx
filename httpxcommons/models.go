package httpxcommons

type BaseResponse interface {
	IsSuccess() bool
	GetMessage() string
	GetPayload() interface{}
	GetRequestId() string
	GetCode() int64
}

func (a *AjaxResponse) IsSuccess() bool {
	return a.Success
}

func (a *AjaxResponse) GetMessage() string {
	return a.Message
}

func (a *AjaxResponse) GetPayload() interface{} {
	return a.Payload
}

func (a *AjaxResponse) GetRequestId() string {
	return a.RequestID
}

func (a *AjaxResponse) GetCode() int64 {
	return a.Code
}

type AjaxResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Payload   interface{} `json:"payload"`
	RequestID string      `json:"request_id"`
	Code      int64       `json:"code"`
	Redirect  string      `json:"redirect,omitempty"`
}

type ListData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}
