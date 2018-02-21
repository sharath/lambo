package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/sharath/lambo/models/external/CMC"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		f, _ := CMC.FetchEntries(25)
		c.JSON(http.StatusOK, f)
	})
	router.Run()
}
