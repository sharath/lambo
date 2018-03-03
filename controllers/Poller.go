package controllers

import (
	"github.com/sharath/lambo/models/extern/CMC"
	"time"
	"fmt"
)

// Poller is a struct for handling updates on CMC
type Poller struct {
	last   int
	gdata  *CMC.GlobalData
	Update chan int
	Resume chan struct{}
	Pause  chan struct{}
	paused bool
}

// StartPoller initializes the poller and returns it
func NewPoller() *Poller {
	p := new(Poller)
	p.gdata = new(CMC.GlobalData)
	p.Update = make(chan int)
	p.Resume = make(chan struct{})
	p.Pause = make(chan struct{})
	p.gdata = CMC.FetchStats()
	p.last = p.gdata.LastUpdated
	return p
}


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
		select {
		case <-p.Resume:
			fmt.Println("Resuming")
			p.paused = false
		case <-p.Pause:
			fmt.Println("Pausing")
			p.paused = true
		}
	}

	go listen()
	go start()
}
