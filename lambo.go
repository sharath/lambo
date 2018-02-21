package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/sharath/lambo/models/external/CMC"
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

	for range controllers.NewPoller().Update {
		fmt.Println("found one lul")
	}

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		f, _ := CMC.FetchEntries(25)
		c.JSON(http.StatusOK, f)
	})
	router.Run()
}
