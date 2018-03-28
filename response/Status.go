package response

import (
	"time"
)

// Status is the response format for the current status of the server
type Status struct {
	Message string `json:"message"`
	Time    int64  `json:"time"`
}

// NewStatus returns a filled status object
func NewStatus(status string) *Status {
	s := new(Status)
	s.Time = time.Now().Unix()
	s.Message = status
	return s
}
