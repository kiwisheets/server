package server

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/kiwisheets/util"
	"gorm.io/gorm"
)

var router *gin.Engine
var port string
var endpoint string

func Setup(gqlHandler *handler.Server, cfg *util.GqlConfig, db *gorm.DB) *gin.RouterGroup {
	router = gin.Default()
	port = cfg.Port
	endpoint = cfg.APIPath

	registerMiddleware(&router.RouterGroup, db, cfg)

	registerRoutes(gqlHandler, &router.RouterGroup, cfg, db)

	return &router.RouterGroup
}

// Run starts a new server
func Run() {
	log.Println("Server listening @ \"/" + endpoint + "\" on " + port)
	router.Run(":" + port)
}
