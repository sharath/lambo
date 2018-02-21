package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sharath/lambo/controllers"
	"github.com/sharath/lambo/models/intern"
	"github.com/sharath/lambo/util"
	"gopkg.in/mgo.v2"
	"net/http"
)

var database *mgo.Database

func main() {
	s, err := mgo.Dial("localhost")
	defer s.Close()
	if err != nil {
		util.HandleError(err, true)
	}
	database = s.DB("lambo")
	lim := 25

	go controllers.StartMongoUpdater(database, lim)

	router := gin.Default()
	router.LoadHTMLGlob("views/templates/*")
	router.Static("/static", "views/static")
	router.GET("/", func(c *gin.Context) {
		var me intern.MongoEntry
		database.C("entries").Find(nil).One(&me)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"TotalMarketCapUsd":            me.Global.TotalMarketCapUsd,
			"Total24HVolumeUsd":            me.Global.Total24HVolumeUsd,
			"BitcoinPercentageOfMarketCap": me.Global.BitcoinPercentageOfMarketCap,
			"ActiveCurrencies":             me.Global.ActiveCurrencies,
			"ActiveAssets":                 me.Global.ActiveAssets,
			"ActiveMarkets":                me.Global.ActiveMarkets,
			"LastUpdated":                  me.Global.LastUpdated,
		})
	})
	router.Run()
}
