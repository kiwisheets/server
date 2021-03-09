package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kiwisheets/util"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Server struct {
	router     *gin.Engine
	config     *util.GqlConfig
	db         *gorm.DB
	gqlHandler *handler.Server
}

func Setup(gqlHandler *handler.Server, cfg *util.GqlConfig, db *gorm.DB) *Server {
	// disable unnecessary debug logging from gin
	gin.SetMode(gin.ReleaseMode)

	var s Server
	s.router = gin.Default()
	s.config = cfg
	s.db = db
	s.gqlHandler = gqlHandler

	// register cors middleware for Apollo Studio if in Dev
	if cfg.Environment == "development" {
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{
			"https://studio.apollographql.com",
		}
		config.AllowCredentials = true
		config.AllowHeaders = append(config.AllowHeaders, "user")
		s.router.Use(cors.New(config))
	}

	registerMiddleware(&s.router.RouterGroup, db, cfg)

	return &s
}

func (s *Server) RegisterMiddleware(middleware ...gin.HandlerFunc) {
	s.router.RouterGroup.Use(middleware...)
}

// Run starts a new server
func (s *Server) Run(log *logrus.Entry) {
	registerRoutes(s.gqlHandler, &s.router.RouterGroup, s.config, s.db)

	SetHealthStatus(HealthStarting)

	log.Println("Server listening @ \"" + s.config.APIPath + "\" on " + s.config.Port)
	s.router.Run(":" + s.config.Port)

	SetHealthStatus(HealthHealthy)
}
