package server

import (
	"github.com/gin-gonic/gin"
	"github.com/kiwisheets/auth"
	"github.com/kiwisheets/util"
	"gorm.io/gorm"
)

func registerMiddleware(router *gin.RouterGroup, db *gorm.DB, cfg *util.GqlConfig) {
	router.Use(auth.Middleware())
}
