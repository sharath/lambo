package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sharath/lambo/controllers"
	"github.com/sharath/lambo/models/intern"
	"github.com/sharath/lambo/util"
	"gopkg.in/mgo.v2"
	"net/http"
	"fmt"
)

var database *mgo.Database
var initalRun bool

func main() {
	s, err := mgo.Dial("localhost")
	defer s.Close()
	if err != nil {
		util.HandleError(err, true)
	}
	database = s.DB("lambo")
	count, _ := database.C("users").Count()
	if count < 1 {
		initalRun = true
	}
	lim := 25

	go controllers.StartMongoUpdater(database, lim)

	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	router := gin.Default()
	router.LoadHTMLGlob("views/templates/*")
	router.Static("/static", "views/static")
	router.GET("/", login)
	router.POST("/authenticate", authenticate)
	router.POST("/register", register)
	router.Run(":80")
}

func login(context *gin.Context) {
	if !initalRun {
		context.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title": "Login",
		})
	} else {
		context.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title": "Register",
		})
	}
}

func register(context *gin.Context) {
	if !initalRun {
		context.JSON(http.StatusUnauthorized, util.NewUnauthorizedResponse())
		return
	}
	u := context.PostForm("username")
	p := context.PostForm("password")
	user, err := intern.CreateUser(u, p, database.C("users"))
	if err != nil {
		context.JSON(http.StatusUnauthorized, util.NewUnauthorizedResponse())
		return
	}
	context.JSON(http.StatusAccepted, user)
}

func authenticate(context *gin.Context) {
	u := context.PostForm("username")
	p := context.PostForm("password")
	authKey, err := intern.AuthenticateUser(u, p, database.C("users"))
	if err != nil {
		context.JSON(http.StatusUnauthorized, util.NewUnauthorizedResponse())
		return
	}
	context.JSON(http.StatusAccepted, gin.H{
		"key": authKey,
	})

	cookie1 := &http.Cookie{Name: "sample", Value: "sample", HttpOnly: false}
	http.SetCookie(context.Writer, cookie1)
	fmt.Println(context.Request.Cookie("sample"))
}

func dashboard(context *gin.Context) {
	u, err := context.Request.Cookie("auth")
	if err != nil {
		login(context)
		return
	}
	fmt.Println(u.Value)
}
