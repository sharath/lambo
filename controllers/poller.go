package controllers

import "time"

type Poller struct {
	last string
}

func NewPoller() *Poller {
	return new(Poller)
}

func Poll() {
	poll := func() {
		for range time.NewTicker(time.Second).C {
			print(time.Now())
		}
	}
	go poll()
}