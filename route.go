package server

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/kiwisheets/util"
	"gorm.io/gorm"
)

func graphqlHandler(gqlHandler *handler.Server, db *gorm.DB, cfg *util.GqlConfig) gin.HandlerFunc {
	gql := GraphqlHandler(gqlHandler, db, cfg)
	return func(c *gin.Context) {
		gql.ServeHTTP(c.Writer, c.Request)
	}
}

func healthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

func playgroundHandler(cfg *util.GqlConfig) gin.HandlerFunc {
	playground := playground.Handler("GraphQL playground", cfg.PlaygroundAPIPath)
	return func(c *gin.Context) {
		playground.ServeHTTP(c.Writer, c.Request)
	}
}

func registerRoutes(gqlHandler *handler.Server, router *gin.RouterGroup, cfg *util.GqlConfig, db *gorm.DB) {
	router.GET("/health", healthHandler())

	// support GET for automatic persisted queries
	router.GET(cfg.APIPath, graphqlHandler(gqlHandler, db, cfg))
	router.POST(cfg.APIPath, graphqlHandler(gqlHandler, db, cfg))

	if cfg.PlaygroundEnabled {
		router.GET(cfg.PlaygroundPath, playgroundHandler(cfg))
	}
}
