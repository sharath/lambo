package controllers

import (
	"encoding/json"
	"github.com/sharath/lambo/models/extern/CMC"
	"github.com/sharath/lambo/models/intern"
	"gopkg.in/mgo.v2"
)

// MongoUpdater adds MongoEntries into MongoDB on each tick from Poller
type MongoUpdater struct {
	db  *mgo.Database
	StartSignal chan struct {}
	PauseSignal chan struct {}
	lim int
}

// StartMongoUpdater initializes the MongoUpdater and returns it
func StartMongoUpdater(db *mgo.Database, lim int) *MongoUpdater {
	m := new(MongoUpdater)
	m.db = db
	m.lim = lim
	go m.start()
	return m
}

func (m *MongoUpdater) start() {
	// every time there's an update from poller
	var me intern.MongoEntry
	me.Tokens = make([]*intern.Token, m.lim)
	me.Global = new(intern.Global)
	var global *CMC.GlobalData
	var entries []*CMC.Entry
	for range StartPoller().Update {
		// get values
		entries = CMC.FetchEntries(m.lim)
		global = CMC.FetchStats()

		// set intern values to extern ones
		t, _ := json.Marshal(global)
		json.Unmarshal(t, &me.Global)

		t, _ = json.Marshal(entries)
		json.Unmarshal(t, &me.Tokens)

		m.db.C("entries").Insert(me)
	}
}
