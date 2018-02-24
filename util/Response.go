package util

type Response struct {
	Message string `json:"message"`
}

func NewUnauthorizedResponse() *Response {
	r := new(Response)
	r.Message = "Unauthorized"
	return r
}