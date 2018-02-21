package controllers

import (
	"gopkg.in/mgo.v2"
	"github.com/sharath/lambo/models/extern/CMC"
	"github.com/sharath/lambo/models/intern"
	"encoding/json"
	"fmt"
)

type MongoUpdater struct {
	db  *mgo.Database
	lim int
}

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

	var global *CMC.GlobalData
	var entries []*CMC.Entry
	for range StartPoller().Update {
		// get values
		entries = CMC.FetchEntries(m.lim)
		global.FetchStats()

		// set intern values to extern ones
		t, _ := json.Marshal(global)
		json.Unmarshal(t, &me.Global)

		t, _ = json.Marshal(entries)
		json.Unmarshal(t, &me.Tokens)

		m.db.C("entries").Insert(me)
	}
}
