package server

import (
	"github.com/gin-gonic/gin"
	"github.com/kiwisheets/auth"
)

func registerMiddleware(router *gin.RouterGroup) {
	router.Use(auth.Middleware())
}
