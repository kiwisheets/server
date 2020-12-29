package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kiwisheets/auth"
	"github.com/kiwisheets/util"
	"gorm.io/gorm"
)

func registerMiddleware(router *gin.RouterGroup, db *gorm.DB, cfg *util.GqlConfig) {
	router.Use(auth.Middleware())

	// register cors middleware for Apollo Studio if in Dev
	if cfg.Environment == "development" {
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{
			"https://studio.apollographql.com",
		}
		config.AllowCredentials = true
		router.Use(cors.New(config))
	}
}
