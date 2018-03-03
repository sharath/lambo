package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sharath/lambo/controllers"
	"github.com/sharath/lambo/models/intern"
	"github.com/sharath/lambo/util"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
	"github.com/sharath/lambo/binders"
)

var database *mgo.Database
var updater *controllers.MongoUpdater

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
		port = ":8080"
	} else {
		gin.SetMode(gin.ReleaseMode)
		port = ":80"
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

func initialRun() bool {
	count, _ := database.C("users").Count()
	if count != 0 {
		return false
	}
	return true
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
func register(context *gin.Context) {
	if initialRun() {
		u := context.PostForm("username")
		p := context.PostForm("password")
		_, err := intern.CreateUser(u, p, database.C("users"))
		if err == nil {
			context.HTML(http.StatusOK, "login.tmpl", binders.GetLoginBinding(false, true))
			return
		}
		context.HTML(http.StatusOK, "register.tmpl", binders.GetRegisterBinding(false))
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
func forceLogin(context *gin.Context) {
	context.HTML(http.StatusUnauthorized, "login.tmpl", binders.GetLoginBinding(true, false))
}

func login(context *gin.Context) {
	if initialRun() {
		context.HTML(http.StatusOK, "register.tmpl", binders.GetRegisterBinding(false))
		return
	}
	context.HTML(http.StatusOK, "login.tmpl", binders.GetLoginBinding(false, false))
}
func dashboard(context *gin.Context) {
	if !validsession(context) {
		forceLogin(context)
		return
	}
	context.HTML(http.StatusOK, "dashboard.tmpl", binders.GetDashboardBinding(database, updater))
}
