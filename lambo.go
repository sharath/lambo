package main

import (
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/sharath/lambo/controllers"
	"github.com/sharath/lambo/models/intern"
	"github.com/sharath/lambo/util"
	"gopkg.in/mgo.v2"
	"net/http"
	"time"
	"os"
)

var database *mgo.Database
var updater *controllers.MongoUpdater

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
	updater = controllers.NewMongoUpdater(database, 25)
	updater.Start()

	prod := os.Getenv("LAMBO_PROD")
	var port string
	if prod != "" {
		gin.SetMode(gin.DebugMode)
		port = "8080"
	} else {
		gin.SetMode(gin.ReleaseMode)
		port = "80"
	}

	router := gin.Default()
	router.LoadHTMLGlob("views/templates/*")
	router.Static("/static", "views/static")

	router.GET("/", login)
	router.GET("/dashboard", dashboard)

	router.POST("/authenticate", authenticate)
	router.POST("/register", register)
	router.POST("/dashboard/", signal)

	router.Run(port)
}

func signal(context *gin.Context) {
	if !validsession(context) {
		forceLogin(context)
		return
	}
	action := context.PostForm("action")
	if action == "pause" {
		updater.Pause()
	} else if action == "resume" {
		updater.Resume()
	}
	dashboard(context)
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
		forceLogin(context)
		return
	}
	auth := &http.Cookie{Name: "auth", Value: authKey, HttpOnly: false}
	user := &http.Cookie{Name: "id", Value: id, HttpOnly: false}
	http.SetCookie(context.Writer, auth)
	http.SetCookie(context.Writer, user)
	context.Redirect(303, "dashboard")
}

func validsession(context *gin.Context) bool {
	authKey, err := context.Cookie("auth")
	if err != nil {
		return false
	}
	id, err := context.Cookie("id")
	if err != nil {
		return false
	}
	valid, err := intern.VerifyAuthKey(id, authKey, database.C("users"))
	if err != nil || !valid {
		return false
	}
	return true
}

func dashboard(context *gin.Context) {
	if !validsession(context) {
		forceLogin(context)
		return
	}

	// TODO make the dashboard template look nicer, add control buttons and endpoints
	var e intern.MongoEntry
	database.C("entries").Find(nil).Sort("-global_data.last_updated").One(&e)
	var g intern.Global
	if e.Global != nil {
		g = *e.Global
	} else {
		context.Writer.Write([]byte("Please wait until first entry is fetched. Usually takes about 5 minutes."))
		return
	}
	NumberOfEntries, _ := database.C("entries").Count()
	update := time.Unix(int64(g.LastUpdated), 0)
	context.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
		"TotalMarketCapUsd":            humanize.Commaf(g.TotalMarketCapUsd),
		"Total24HVolumeUsd":            humanize.Commaf(g.Total24HVolumeUsd),
		"BitcoinPercentageOfMarketCap": g.BitcoinPercentageOfMarketCap,
		"ActiveCurrencies":             g.ActiveCurrencies,
		"ActiveAssets":                 g.ActiveAssets,
		"ActiveMarkets":                g.ActiveMarkets,
		"LastUpdated":                  humanize.Time(update),
		"NumberOfEntries":              NumberOfEntries,
		"MongoUpdateStatus":            updater.Status(),
	})
}

func forceLogin(context *gin.Context) {
	context.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
		"unauthorized":   true,
		"justregistered": false,
	})
}
