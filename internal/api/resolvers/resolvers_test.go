package resolvers

// parseTestSchema parse schema and creates the root resolvers with needed dependencies.
/* func parseTestSchema(db *sql.DB) *graphql.Schema {
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
} */
