package testutils

import (
	"testing"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/api/resolvers"
	"github.com/brice-74/golang-base-api/internal/api/schema"
	"github.com/google/go-cmp/cmp"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/errors"
)

func ParseTestSchema(app *application.Application) *graphql.Schema {
	return graphql.MustParseSchema(
		schema.String(),
		&resolvers.Root{
			App: app,
		},
	)
}

func TestGqlError(t *testing.T, qerr *errors.QueryError, expect *ExpectResolverError) {
	if qerr != nil {
		if qerr.Message != expect.Msg {
			t.Errorf("expect query error message: %s, got: %s", expect.Msg, qerr.Message)
		}
		if expect.Extensions != nil {
			if diff := cmp.Diff(expect.Extensions, qerr.Extensions); diff != "" {
				t.Error(diff)
			}
		}
	} else {
		t.Fatal("resolver error expected but not found")
	}
}

type ExpectResolverError struct {
	Msg        string
	Extensions map[string]interface{}
}
