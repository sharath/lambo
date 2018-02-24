package main

import (
	"fmt"
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
	router.GET("/", login)
	router.GET("/register", getRegister)
	router.POST("/authenticate", authenticate)
	router.POST("/register", register)
	router.Run()
}

func login(context *gin.Context) {
	context.HTML(http.StatusOK, "login.tmpl", gin.H{
		"title": "Login",
	})
}

func getRegister(context *gin.Context) {
	// TODO
}

func register(context *gin.Context) {
	// TODO
}

func authenticate(context *gin.Context) {
	u := context.PostForm("username")
	p := context.PostForm("password")
	fmt.Println(u, p)
	authKey, err := intern.AuthenticateUser(u, p, database.C("users"))
	if err != nil {
		context.JSON(http.StatusUnauthorized, util.NewUnauthorizedResponse())
		return
	}
	context.JSON(http.StatusAccepted, gin.H{
		"key": authKey,
	})
}
