package controllers

import (
	"time"
	"github.com/sharath/lambo/models/external/CMC"
)

type Poller struct {
	last   int
	gdata  *CMC.GlobalData
	Update chan int
}

func NewPoller() *Poller {
	p := new(Poller)
	p.gdata = new(CMC.GlobalData)
	p.Update = make(chan int)
	p.gdata.Update()
	p.last = p.gdata.LastUpdated
	go p.start()
	return p
}

func (p *Poller) start() {
	for range time.NewTicker(time.Second * 10).C {
		p.gdata.Update()
		if p.gdata.LastUpdated != p.last {
			p.last = p.gdata.LastUpdated
			p.Update <- p.last
		}
	}
}
