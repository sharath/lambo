package util

// Response contains a message
type Response struct {
	Message string `json:"message"`
}

// NewUnauthorizedResponse returns a response that says "Unauthorized"
func NewUnauthorizedResponse() *Response {
	r := new(Response)
	r.Message = "Unauthorized"
	return r
}
