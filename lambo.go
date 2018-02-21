package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sharath/lambo/controllers"
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
	router.GET("/", login)
	router.POST("/authenticate", authenticate)
	router.Run()
}

func login(context *gin.Context) {
	context.HTML(http.StatusOK, "login.tmpl", gin.H{
		"title": "Login",
	})
}

func authenticate(context *gin.Context) {
	u := gin.H{
		"username": context.PostForm("username"),
		"password": context.PostForm("password"),
	}
	context.JSON(http.StatusUnauthorized, u)
}
