package server

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/kiwisheets/util"
)

// GraphqlHandler constructs and returns a http handler
func GraphqlHandler(gqlHandler *handler.Server, cfg *util.GqlConfig) http.Handler {
	gqlHandler.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	gqlHandler.AddTransport(transport.Options{})
	gqlHandler.AddTransport(transport.GET{})
	gqlHandler.AddTransport(transport.POST{})
	gqlHandler.AddTransport(transport.MultipartForm{})

	gqlHandler.Use(extension.AutomaticPersistedQuery{
		Cache: cfg.Cache,
	})

	gqlHandler.Use(extension.Introspection{})
	gqlHandler.Use(&extension.ComplexityLimit{
		Func: func(ctx context.Context, rc *graphql.OperationContext) int {
			return cfg.ComplexityLimit
		},
	})

	return gqlHandler
}
