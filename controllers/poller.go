package controllers

import (
	"time"
	"fmt"
)

type Poller struct {
	last string
}

func NewPoller() *Poller {
	return new(Poller)
}

func (p *Poller) Poll() {
	poll := func() {
		for range time.NewTicker(time.Second).C {
			fmt.Println(time.Now())
		}
	}
	go poll()
}