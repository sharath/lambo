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

func initialRun() bool {
	count, _ := database.C("users").Count()
	if count != 0 {
		return false
	}
	return true
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

	gin.SetMode(gin.ReleaseMode)
	//gin.SetMode(gin.DebugMode)

	router := gin.Default()
	router.LoadHTMLGlob("views/templates/*")
	router.Static("/static", "views/static")
	router.GET("/", login)
	router.POST("/authenticate", authenticate)
	router.POST("/register", register)
	router.GET("/dashboard", dashboard)
	router.Run(":8080")
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
	id, authKey, err := intern.AuthenticateUser(u, p, database.C("users"))
	if err != nil {
		forcelogin(context)
		return
	}
	auth := &http.Cookie{Name: "auth", Value: authKey, HttpOnly: false}
	user := &http.Cookie{Name: "id", Value: id, HttpOnly: false}
	http.SetCookie(context.Writer, auth)
	http.SetCookie(context.Writer, user)
	context.Redirect(303, "dashboard")
}

func dashboard(context *gin.Context) {
	authKey, err := context.Cookie("auth")
	if err != nil {
		forcelogin(context)
		return
	}
	id, err := context.Cookie("id")
	if err != nil {
		forcelogin(context)
		return
	}

	valid, err := intern.VerifyAuthKey(id, authKey, database.C("users"))
	if err != nil || !valid {
		forcelogin(context)
		return
	}

	// TODO make the dashboard template look nicer, add control buttons and endpoints
	var e intern.MongoEntry
	database.C("entries").Find(nil).Sort("-global_data.last_updated").One(&e)
	var g intern.Global
	if e.Global != nil {
		g = *e.Global
	} else {
		context.Writer.Write([]byte("please wait until first entry is fetched"))
		return
	}
	context.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
		"TotalMarketCapUsd":            g.TotalMarketCapUsd,
		"Total24HVolumeUsd":            g.Total24HVolumeUsd,
		"BitcoinPercentageOfMarketCap": g.BitcoinPercentageOfMarketCap,
		"ActiveCurrencies":             g.ActiveCurrencies,
		"ActiveAssets":                 g.ActiveAssets,
		"ActiveMarkets":                g.ActiveMarkets,
		"LastUpdated":                  g.LastUpdated,
	})
}

func forcelogin(context *gin.Context) {
	context.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
		"unauthorized":   true,
		"justregistered": false,
	})
}
