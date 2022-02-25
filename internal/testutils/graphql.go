package testutils

import (
	"database/sql"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/api/resolvers"
	"github.com/brice-74/golang-base-api/internal/api/schema"
	"github.com/graph-gophers/graphql-go"
)

func ParseTestSchema(db *sql.DB) *graphql.Schema {
	app := &application.Application{
		Models: application.NewModels(db),
	}

	return graphql.MustParseSchema(
		schema.String(),
		&resolvers.Root{
			App: app,
		},
	)
}

/* func ParseTestSchemaWithCustomRoot(root *resolvers.Root) *graphql.Schema {
	return graphql.MustParseSchema(
		schema.String(),
		root,
	)
} */
