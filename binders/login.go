package binders

import (
	"github.com/gin-gonic/gin"
)

func GetLoginBinding(unauthorized, new bool) gin.H {
	return gin.H{
		"unauthorized":   false,
		"justregistered": false,
	}
}
