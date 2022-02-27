package testutils

import (
	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/api/resolvers"
	"github.com/brice-74/golang-base-api/internal/api/schema"
	"github.com/graph-gophers/graphql-go"
)

func ParseTestSchema(app *application.Application) *graphql.Schema {
	return graphql.MustParseSchema(
		schema.String(),
		&resolvers.Root{
			App: app,
		},
	)
}
