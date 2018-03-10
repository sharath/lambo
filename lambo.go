package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sharath/lambo/database"
	"github.com/sharath/lambo/poller"
	"github.com/sharath/lambo/response"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
)

var lambo *mgo.Database
var users *mgo.Collection
var updater *poller.MongoUpdater

func main() {
	s, err := mgo.Dial("localhost")
	defer s.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	lambo = s.DB("lambo")
	users = lambo.C("users")
	updater = poller.NewMongoUpdater(lambo, 25)
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

	r := gin.Default()
	r.GET("/", root)
	r.POST("/register", register)
	r.POST("/login", login)
	r.Run(port)
}

func root(c *gin.Context) {
	if lambo.Session.Ping() != nil {
		c.JSON(http.StatusInternalServerError, response.NewStatus("Error Encountered"))
		return
	}
	c.JSON(http.StatusOK, response.NewStatus("Normal Operation"))
}

func register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username != "" && password != "" {
		u, err := database.CreateUser(username, password, users)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.NewStatus(err.Error()))
			return
		}
		authKey := u.Login(password, users)
		c.JSON(http.StatusOK, response.NewLogin(authKey))
		return
	}
	c.JSON(http.StatusBadRequest, response.NewStatus("missing username or password"))
	return
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username != "" && password != "" {
		u := database.FetchUser(username, users)
		if u.Username != "" {
			authKey := u.Login(password, users)
			if authKey != "" {
				c.JSON(http.StatusOK, response.NewLogin(authKey))
				return
			}
		}
	}
	c.JSON(http.StatusBadRequest, response.NewStatus("invalid credentials"))
	return
}
