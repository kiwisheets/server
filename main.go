package server

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/kiwisheets/util"
	"gorm.io/gorm"
)

var router *gin.Engine

func Setup(gqlHandler *handler.Server, cfg *util.GqlConfig, db *gorm.DB) *gin.RouterGroup {
	router = gin.Default()

	registerMiddleware(&router.RouterGroup, db, cfg)

	registerRoutes(gqlHandler, &router.RouterGroup, cfg, db)

	return &router.RouterGroup
}

// Run starts a new server
func Run(cfg *util.GqlConfig, db *gorm.DB) {
	log.Println("Server listening @ \"/\" on " + cfg.Port)
	router.Run()
}
