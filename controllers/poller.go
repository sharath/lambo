package controllers

import (
	"time"
	"github.com/sharath/lambo/models/external/CMC"
)

type Poller struct {
	last int
	gdata *CMC.GlobalData
}

func NewPoller() *Poller {
	p := new(Poller)
	p.gdata.Update()
	p.last = p.gdata.LastUpdated
	return p
}

func (p *Poller) GetNextUpdate() bool {
	updated := make(chan bool)
	poll := func() {
		for range time.NewTicker(time.Minute).C {
			p.gdata.Update()
			if p.gdata.LastUpdated != p.last {
				p.last = p.gdata.LastUpdated
				updated <- true
			}
		}
	}
	go poll()
	return <- updated
}