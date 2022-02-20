package api

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/api/resolvers"
	"github.com/brice-74/golang-base-api/internal/api/schema"
)

// GraphQL is the main entrypoint for queries and mutations.
func GraphQL(app *application.Application) http.HandlerFunc {
	opts := []graphql.SchemaOpt{graphql.Logger(logger{app: app})}

	s := graphql.MustParseSchema(
		schema.String(),
		&resolvers.Root{
			App: app,
		},
		opts...,
	)

	return func(w http.ResponseWriter, r *http.Request) {
		h := relay.Handler{Schema: s}
		h.ServeHTTP(w, r)
	}
}

// logger for GraphQL
type logger struct {
	app *application.Application
}

// LogPanic is used to log recovered panic values that occur during query execution.
func (l logger) LogPanic(ctx context.Context, value interface{}) {
	const size = 64 << 10
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]
	err := fmt.Errorf("graphql: panic occurred: %v\n%s\ncontext: %v", value, buf, ctx)
	l.app.Logger.PrintError(err, nil)
}
