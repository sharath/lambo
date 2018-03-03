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
	Resume chan byte
	Pause  chan byte
	paused bool
}

// NewPoller initializes the poller and returns it
func NewPoller() *Poller {
	p := new(Poller)
	p.paused = true
	p.gdata = new(CMC.GlobalData)
	p.Update = make(chan int)
	p.Resume = make(chan byte)
	p.Pause = make(chan byte)
	p.gdata = CMC.FetchStats()
	p.last = p.gdata.LastUpdated
	return p
}

// Starts the Poller and listens for status updates
func (p *Poller) Start() {
	start := func() {
		for range time.NewTicker(time.Second * 10).C {
			if !p.paused {
				p.gdata = CMC.FetchStats()
				if p.gdata.LastUpdated != p.last {
					p.last = p.gdata.LastUpdated
					p.Update <- p.last
				}
			}
		}
	}
	listen := func() {
		for {
			select {
			case <-p.Resume:
				p.paused = false
			case <-p.Pause:
				p.paused = true
			}
		}
	}
	go listen()
	go start()
}
