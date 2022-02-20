package schema

import "embed"

//go:embed schema.graphql
var schemaFS embed.FS

func String() string {
	f, err := schemaFS.ReadFile("schema.graphql")
	if err != nil {
		panic(err)
	}

	return string(f)
}
