package controllers

import (
	"gopkg.in/mgo.v2"
	"github.com/sharath/lambo/models/extern/CMC"
	"github.com/sharath/lambo/models/intern"
)

type MongoUpdater struct {
	db *mgo.Database
}

func StartMongoUpdater(db *mgo.Database, lim int) *MongoUpdater {
	m := new(MongoUpdater)
	m.db = db
	go m.start()
	return m
}

func (m *MongoUpdater) start() {
	// every time there's an update from poller
	for range NewPoller().Update {
		var me intern.MongoEntry
		var glob *CMC.GlobalData
		glob.Update()
		me.Global.TotalMarketCapUsd = glob.TotalMarketCapUsd
		me.Global.Total24HVolumeUsd = glob.Total24HVolumeUsd
		me.Global.BitcoinPercentageOfMarketCap = glob.BitcoinPercentageOfMarketCap
		me.Global.ActiveCurrencies = glob.ActiveCurrencies
		me.Global.ActiveAssets = glob.ActiveAssets
		me.Global.ActiveMarkets = glob.ActiveMarkets
		me.Global.LastUpdated = glob.LastUpdated
	}
}
