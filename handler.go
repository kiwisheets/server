package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/kiwisheets/util"
	"gorm.io/gorm"
)

// GraphqlHandler constructs and returns a http handler
func GraphqlHandler(gqlHandler *handler.Server, db *gorm.DB, cfg *util.GqlConfig) http.Handler {
	// init APQ cache
	cache, err := newCache(cfg.Redis.Address, 24*time.Hour)
	if err != nil {
		panic(fmt.Errorf("cannot create APQ cache: %v", err))
	}

	gqlHandler.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	gqlHandler.AddTransport(transport.Options{})
	gqlHandler.AddTransport(transport.GET{})
	gqlHandler.AddTransport(transport.POST{})
	gqlHandler.AddTransport(transport.MultipartForm{})

	gqlHandler.Use(extension.AutomaticPersistedQuery{
		Cache: cache,
	})

	gqlHandler.Use(extension.Introspection{})
	gqlHandler.Use(&extension.ComplexityLimit{
		Func: func(ctx context.Context, rc *graphql.OperationContext) int {
			return cfg.ComplexityLimit
		},
	})

	return gqlHandler
}
