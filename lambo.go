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

	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	router := gin.Default()
	router.LoadHTMLGlob("views/templates/*")
	router.Static("/static", "views/static")
	router.GET("/", login)
	router.POST("/authenticate", authenticate)
	router.POST("/register", register)
	router.GET("/dashboard", dashboard)
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

	cookie := &http.Cookie{Name: "auth", Value: authKey, HttpOnly: false, Secure: true}
	http.SetCookie(context.Writer, cookie)

	dashboard(context)
}

func dashboard(context *gin.Context) {
	// TODO make the dashboard template look nicer, add control buttons and endpoints
	var e intern.MongoEntry
	database.C("entries").Find(nil).One(&e)
	var g intern.Global
	if e.Global != nil {
		g = *e.Global
	} else {
		// to prevent panics when mongodb is empty
		// TODO do this better
		g = intern.Global{
			TotalMarketCapUsd:            0,
			Total24HVolumeUsd:            0,
			BitcoinPercentageOfMarketCap: 0,
			ActiveCurrencies:             0,
			ActiveAssets:                 0,
			ActiveMarkets:                0,
			LastUpdated:                  0,
		}
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
