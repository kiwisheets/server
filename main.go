package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kiwisheets/util"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var router *gin.Engine
var port string
var endpoint string

func Setup(gqlHandler *handler.Server, cfg *util.GqlConfig, db *gorm.DB) *gin.RouterGroup {
	// disable unnecessary debug logging from gin
	gin.SetMode(gin.ReleaseMode)

	router = gin.Default()
	port = cfg.Port
	endpoint = cfg.APIPath

	registerMiddleware(&router.RouterGroup, db, cfg)

	// register cors middleware for Apollo Studio if in Dev
	if cfg.Environment == "development" {
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{
			"https://studio.apollographql.com",
		}
		config.AllowCredentials = true
		config.AllowHeaders = append(config.AllowHeaders, "user")
		router.Use(cors.New(config))
	}

	registerRoutes(gqlHandler, &router.RouterGroup, cfg, db)

	return &router.RouterGroup
}

// Run starts a new server
func Run(log *logrus.Entry) {
	SetHealthStatus(HealthStarting)

	log.Println("Server listening @ \"" + endpoint + "\" on " + port)
	router.Run(":" + port)

	SetHealthStatus(HealthHealthy)
}
