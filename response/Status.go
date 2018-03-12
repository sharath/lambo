package response

import (
	"time"
)

type Status struct {
	Message string `json:"message"`
	Time    int64  `json:"time"`
}

func NewStatus(status string) *Status {
	s := new(Status)
	s.Time = time.Now().Unix()
	s.Message = status
	return s
}
