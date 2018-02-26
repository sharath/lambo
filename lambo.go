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
var checked bool
func initialRun() bool {
	if !checked {
		count, _ := database.C("users").Count()
		if count != 0 {
			return false
		} else {
			checked = true
			return checked
		}
	}
	return checked
}

func main() {
	s, err := mgo.Dial("localhost")
	defer s.Close()
	if err != nil {
		util.HandleError(err, true)
	}
	database = s.DB("lambo")
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
	if initialRun() {
		context.HTML(http.StatusOK, "register.tmpl", gin.H{
			"invalid": false,
		})
		return
	}
	context.HTML(http.StatusOK, "login.tmpl", gin.H{
		"unauthorized":   false,
		"justregistered": false,
	})
}

func register(context *gin.Context) {
	if initialRun() {
		u := context.PostForm("username")
		p := context.PostForm("password")
		_, err := intern.CreateUser(u, p, database.C("users"))
		if err == nil {
			context.HTML(http.StatusOK, "login.tmpl", gin.H{
				"unauthorized":   false,
				"justregistered": true,
			})
			return
		}
		context.HTML(http.StatusOK, "register.tmpl", gin.H{
			"invalid": true,
		})
		return
	}
}

func authenticate(context *gin.Context) {
	u := context.PostForm("username")
	p := context.PostForm("password")
	authKey, err := intern.AuthenticateUser(u, p, database.C("users"))
	if err != nil {
		context.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
			"unauthorized":   true,
			"justregistered": false,
		})
		return
	}
	// TODO set a cookie since authorized
	// the below doesn't work for some reason

	cookie := &http.Cookie{Name: "auth", Value: authKey, HttpOnly: false}
	http.SetCookie(context.Writer, cookie)

	dashboard(context)
}

func dashboard(context *gin.Context) {
	context.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
		"TotalMarketCapUsd": "testing",
		"Total24HVolumeUsd": "testing",
		"BitcoinPercentageOfMarketCap": "testing",
		"ActiveCurrencies": "testing",
		"ActiveAssets": "testing",
		"ActiveMarkets": "testing",
		"LastUpdated": "testing",
	})
}
