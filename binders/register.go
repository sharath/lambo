package binders

import "github.com/gin-gonic/gin"

func GetRegisterBinding(invalid bool) gin.H {
	return gin.H{
		"invalid": invalid,
	}
}
