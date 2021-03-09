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
	router   *gin.Engine
	port     string
	endpoint string
}

func Setup(gqlHandler *handler.Server, cfg *util.GqlConfig, db *gorm.DB) *Server {
	// disable unnecessary debug logging from gin
	gin.SetMode(gin.ReleaseMode)

	var s Server

	s.router = gin.Default()
	s.port = cfg.Port
	s.endpoint = cfg.APIPath

	registerMiddleware(&s.router.RouterGroup, db, cfg)

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

	registerRoutes(gqlHandler, &s.router.RouterGroup, cfg, db)

	return &s
}

func (s *Server) RegisterMiddleware(middleware ...gin.HandlerFunc) {
	s.router.RouterGroup.Use(middleware...)
}

// Run starts a new server
func (s *Server) Run(log *logrus.Entry) {
	SetHealthStatus(HealthStarting)

	log.Println("Server listening @ \"" + s.endpoint + "\" on " + s.port)
	s.router.Run(":" + s.port)

	SetHealthStatus(HealthHealthy)
}
