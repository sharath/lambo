package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"github.com/sharath/lambo/util"
	"github.com/sharath/lambo/controllers"
	"net/http"
	"github.com/sharath/lambo/models/intern"
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
		c.HTML(http.StatusOK, "index.tmpl", me.Global)
	})
	router.Run()
}
