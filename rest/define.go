// define
package rest

type RestResponse struct {
	Request_id int64       `json:"request_id"`
	Code       int         `json:"code"`
	Msg        interface{} `json:"msg"`
}

type Error struct {
	Name             string        `json:"name"`
	Debug_id         int64         `json:"debug_id"`
	Message          string        `json:"message"`
	Information_link string        `json:"information_link"`
	Details          []ErrorDetail `json:"details"`
}

type ErrorDetail struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
}
