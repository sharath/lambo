package controllers

import (
	"github.com/sharath/lambo/models/extern/CMC"
	"time"
)

// Poller is a struct for handling updates on CMC
type Poller struct {
	last   int
	gdata  *CMC.GlobalData
	Update chan int
}

// StartPoller initializes the poller and returns it
func StartPoller() *Poller {
	p := new(Poller)
	p.gdata = new(CMC.GlobalData)
	p.Update = make(chan int)
	p.gdata = CMC.FetchStats()
	p.last = p.gdata.LastUpdated
	go p.start()
	return p
}

func (p *Poller) start() {
	// check every 10 seconds
	for range time.NewTicker(time.Second * 10).C {
		p.gdata = CMC.FetchStats()
		if p.gdata.LastUpdated != p.last {
			p.last = p.gdata.LastUpdated
			p.Update <- p.last
		}
	}
}
