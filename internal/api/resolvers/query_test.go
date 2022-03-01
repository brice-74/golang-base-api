package resolvers_test

import (
	"context"
	"testing"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/testutils"
	"github.com/brice-74/golang-base-api/internal/testutils/mocks"
	"github.com/graph-gophers/graphql-go/gqltesting"
)

func TestQueryCheck(t *testing.T) {
	var (
		app    = &application.Application{}
		schema = testutils.ParseTestSchema(app)
	)

	test := &gqltesting.Test{
		Schema: schema,
		Query: `
			query {
				queryCheck
			}`,
		ExpectedResult: `
			{
				"queryCheck": "ok"
			}`,
	}

	gqltesting.RunTest(t, test)
}

func TestQueryPanic(t *testing.T) {
	var (
		app = &application.Application{
			Logger: mocks.NewLogger(),
		}
		schema = testutils.ParseTestSchema(app)
	)

	t.Run("should return panic error", func(t *testing.T) {
		testErr := &gqltesting.Test{
			Schema: schema,
			Query: `
				query {
					queryPanic(panic: true)
				}`,
		}

		result := testErr.Schema.Exec(context.Background(), testErr.Query, "", nil)

		got := result.Errors[0]
		expect := "panic occurred: I panic !!!"
		if got != nil {
			if got.Message != expect {
				t.Fatalf("got query error message: %s, expect: %s", got.Message, expect)
			}
		} else {
			t.Fatal("expect error, got nil")
		}
	})

	t.Run("should return no panic", func(t *testing.T) {
		test := &gqltesting.Test{
			Schema: schema,
			Query: `
				query {
					queryPanic(panic: false)
				}`,
			ExpectedResult: `
				{
					"queryPanic": "No panic"
				}	
			`,
		}

		gqltesting.RunTest(t, test)
	})
}
