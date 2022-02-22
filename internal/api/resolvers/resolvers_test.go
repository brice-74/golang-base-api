package resolvers

import (
	"database/sql"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/api/schema"
	"github.com/graph-gophers/graphql-go"
)

// parseTestSchema parse schema and creates the root resolvers with needed dependencies.
func parseTestSchema(db *sql.DB) *graphql.Schema {
	app := &application.Application{
		Models: application.NewModels(db),
	}

	return graphql.MustParseSchema(
		schema.String(),
		&Root{
			App: app,
		},
	)
}

func parseTestSchemaWithCustomRoot(root *Root) *graphql.Schema {
	return graphql.MustParseSchema(
		schema.String(),
		root,
	)
}
