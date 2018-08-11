package poller

import (
	"encoding/json"
	"github.com/sharath/lambo/CMC"
	"github.com/sharath/lambo/database"
	"github.com/globalsign/mgo"
)

// MongoUpdater adds MongoEntries into MongoDB on each tick from Poller
type MongoUpdater struct {
	db  *mgo.Database
	P   *Poller
	lim int
}

// Status returns the current status of the mongoupdater
func (m *MongoUpdater) Status() string {
	if m.P.paused {
		return "Paused"
	}
	return "Running"
}

// NewMongoUpdater initializes the MongoUpdater and returns it
func NewMongoUpdater(db *mgo.Database, lim int) *MongoUpdater {
	m := new(MongoUpdater)
	m.db = db
	m.lim = lim
	return m
}

// GetLim returns limit
func (m *MongoUpdater) GetLim() int {
	return m.lim
}

// Start starts the MongoUpdater
func (m *MongoUpdater) Start() {
	// every time there's an update from poller
	start := func() {
		var me database.MongoEntry
		me.Tokens = make([]*database.Token, m.lim)
		me.Global = new(database.Global)
		var global *CMC.GlobalData
		var entries []*CMC.Entry
		m.P = NewPoller()
		m.P.Start()
		for range m.P.Update {
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
	go start()
}

// Resume the MongoUpdater
func (m *MongoUpdater) Resume() {
	m.P.Resume <- '1'
}

// Pause the MongoUpdater
func (m *MongoUpdater) Pause() {
	m.P.Pause <- '1'
}
