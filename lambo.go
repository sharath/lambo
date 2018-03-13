package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	auth "github.com/sharath/lambo/authentication"
	"github.com/sharath/lambo/poller"
	"github.com/sharath/lambo/response"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
	"time"
)

var lambo *mgo.Database
var users *mgo.Collection
var updater *poller.MongoUpdater
var authmatrix auth.Matrix

func main() {
	s, err := mgo.Dial("localhost")
	defer s.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	lambo = s.DB("lambo")
	users = lambo.C("users")
	authmatrix = auth.NewAuthenticationMatrix()
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
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.GET("/", root)
	r.POST("/register", register)
	r.POST("/login", login)
	r.GET("/do/:action", do)
	r.Run(port)
}

func authenticate(c *gin.Context) *auth.User {
	authkey := c.GetHeader("auth_key")
	if authkey == "" {
		return nil
	}
	if user := auth.FindUserByAuthKey(authkey, users, authmatrix); user != nil {
		return user
	}
	return nil
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
		u, err := auth.CreateUser(username, password, users)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.NewStatus(err.Error()))
			return
		}
		authKey := u.Login(password, users)
		c.JSON(http.StatusOK, response.NewLogin(authKey))
		return
	}
	c.JSON(http.StatusBadRequest, response.NewStatus("missing username or password"))
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username != "" && password != "" {
		u := auth.FetchUser(username, users)
		if u.Username != "" {
			authKey := u.Login(password, users)
			if authKey != "" {
				c.JSON(http.StatusOK, response.NewLogin(authKey))
				return
			}
		}
	}
	c.JSON(http.StatusBadRequest, response.NewStatus("invalid credentials"))
}

func do(c *gin.Context) {
	if user := authenticate(c); user != nil {
		action := c.Param("action")
		switch action {
		case "resume":
			updater.Resume()
			time.Sleep(time.Millisecond)
			c.JSON(http.StatusOK, gin.H{"status": updater.Status(), "time": time.Now().Unix()})
			return
		case "pause":
			updater.Pause()
			time.Sleep(time.Millisecond)
			c.JSON(http.StatusOK, gin.H{"status": updater.Status(), "time": time.Now().Unix()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": updater.Status(), "time": time.Now().Unix()})
		return
	}
	c.JSON(http.StatusUnauthorized, response.NewStatus("unauthorized"))
}
