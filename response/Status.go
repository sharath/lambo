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

// Status is the response format for the current status of the server
type RootStatus struct {
	DB     string `json:"database"`
	Poller string `json:"poller"`
	Time   int64  `json:"time"`
}

// NewStatus returns a filled status object
func NewRootStatus(dbstatus, pstatus string) *RootStatus {
	s := new(RootStatus)
	s.Time = time.Now().Unix()
	s.DB = dbstatus
	s.Poller = pstatus
	return s
}
