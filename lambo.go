package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/sharath/lambo/models/extern/CMC"
	"gopkg.in/mgo.v2"
	"github.com/sharath/lambo/util"
	"github.com/sharath/lambo/controllers"
	"fmt"
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
	router.GET("/", func(c *gin.Context) {
		var t []*CMC.Entry
		t = CMC.FetchEntries(lim)
		fmt.Println(t)
		c.JSON(http.StatusOK, t)
	})
	router.Run()
}
