package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kiwisheets/util"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router     *gin.Engine
	config     *util.GqlConfig
	gqlHandler *handler.Server
}

func Setup(gqlHandler *handler.Server, cfg *util.GqlConfig) *Server {
	// disable unnecessary debug logging from gin
	gin.SetMode(gin.ReleaseMode)

	var s Server
	s.router = gin.Default()
	s.config = cfg
	s.gqlHandler = gqlHandler

	// register cors middleware for Apollo Studio if in Dev
	if checkIfDev(cfg.Environment) {
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{
			"https://studio.apollographql.com",
		}
		config.AllowCredentials = true
		config.AllowHeaders = append(config.AllowHeaders, "user")
		s.router.Use(cors.New(config))
	}

	registerMiddleware(&s.router.RouterGroup)

	return &s
}

func (s *Server) RegisterMiddleware(middleware ...gin.HandlerFunc) {
	s.router.RouterGroup.Use(middleware...)
}

// Run starts a new server
func (s *Server) Run(log *logrus.Entry) {
	registerRoutes(s.gqlHandler, &s.router.RouterGroup, s.config)

	SetHealthStatus(HealthStarting)

	log.Println("Server listening @ \"" + s.config.APIPath + "\" on " + s.config.Port)
	s.router.Run(":" + s.config.Port)

	SetHealthStatus(HealthHealthy)
}
