package binders

import (
	"github.com/gin-gonic/gin"
	"github.com/sharath/lambo/models/intern"
	"github.com/dustin/go-humanize"
	"time"
	"gopkg.in/mgo.v2"
	"github.com/sharath/lambo/controllers"
)

func GetDashboardBinding(database *mgo.Database, updater *controllers.MongoUpdater) gin.H {
	var e intern.MongoEntry
	database.C("entries").Find(nil).Sort("-global_data.last_updated").One(&e)
	var g intern.Global
	if e.Global != nil {
		g = *e.Global
	} else {
		return gin.H{
			"TotalMarketCapUsd":            humanize.Commaf(g.TotalMarketCapUsd),
			"Total24HVolumeUsd":            humanize.Commaf(g.Total24HVolumeUsd),
			"BitcoinPercentageOfMarketCap": g.BitcoinPercentageOfMarketCap,
			"ActiveCurrencies":             g.ActiveCurrencies,
			"ActiveAssets":                 g.ActiveAssets,
			"ActiveMarkets":                g.ActiveMarkets,
			"LastUpdated":                  humanize.Time(time.Unix(0, 0)),
			"NumberOfEntries":              0,
			"MongoUpdateStatus":            updater.Status(),
		}
	}
	update := time.Unix(int64(g.LastUpdated), 0)
	NumberOfEntries, _ := database.C("entries").Count()
	return gin.H{
		"TotalMarketCapUsd":            humanize.Commaf(g.TotalMarketCapUsd),
		"Total24HVolumeUsd":            humanize.Commaf(g.Total24HVolumeUsd),
		"BitcoinPercentageOfMarketCap": g.BitcoinPercentageOfMarketCap,
		"ActiveCurrencies":             g.ActiveCurrencies,
		"ActiveAssets":                 g.ActiveAssets,
		"ActiveMarkets":                g.ActiveMarkets,
		"LastUpdated":                  humanize.Time(update),
		"NumberOfEntries":              NumberOfEntries,
		"MongoUpdateStatus":            updater.Status(),
	}
}
