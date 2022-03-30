package resolvers_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/domains/user"
	"github.com/brice-74/golang-base-api/internal/testutils"
	"github.com/brice-74/golang-base-api/internal/testutils/factory"
	"github.com/graph-gophers/graphql-go/gqltesting"
)

func TestMe(t *testing.T) {
	var (
		db     = testutils.PrepareDB(t)
		app    = testutils.NewApplication(db)
		schema = testutils.ParseTestSchema(app)
		fac    = factory.New(t, db)
	)

	u := fac.CreateUserAccount(&user.User{
		Roles: user.Roles{user.RoleUser},
	})

	anonymous := fac.CreateUserAccount(user.AnonymousUser)

	queryString := `
		{
			me {
				id
				createdAt
				updatedAt
				active
				email
				roles
				profilName
				shortId
			}
		}
	`

	tests := []struct {
		title       string
		gqltest     *gqltesting.Test
		expectError *testutils.ExpectResolverError
	}{
		{
			title: "Should return me",
			gqltest: &gqltesting.Test{
				Context: app.ContextWithClient(context.Background(), &application.ClientCtx{User: u}),
				Schema:  schema,
				Query:   queryString,
			},
		},
		{
			title: "Should be unauthorized",
			gqltest: &gqltesting.Test{
				Context: app.ContextWithClient(context.Background(), &application.ClientCtx{User: anonymous}),
				Schema:  schema,
				Query:   queryString,
			},
			expectError: &testutils.ExpectResolverError{
				Msg: "error [Unauthorized]: Unauthorized access",
				Extensions: map[string]interface{}{
					"code":       "Unauthorized",
					"message":    "Unauthorized access",
					"statusCode": 401,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			result := tt.gqltest.Schema.Exec(tt.gqltest.Context, tt.gqltest.Query, "", nil)

			if tt.expectError != nil {
				testutils.TestGqlError(t, result.Errors[0], tt.expectError)
			} else {
				var res MeResponse

				data, _ := result.Data.MarshalJSON()
				if err := json.Unmarshal(data, &res); err != nil {
					t.Fatal(err)
				}

				var uRes = res.Me

				testUserOutput := []struct {
					field string
					err   bool
				}{
					{field: "ID", err: u.ID != uRes.ID},
					{field: "CreatedAt", err: u.CreatedAt.Format(time.RFC3339) != uRes.CreatedAt.Format(time.RFC3339)},
					{field: "UpdatedAt", err: u.UpdatedAt.Format(time.RFC3339) != uRes.UpdatedAt.Format(time.RFC3339)},
					{field: "DeactivatedAt", err: !u.DeactivatedAt.IsZero()},
					{field: "Email", err: u.Email != uRes.Email},
					{field: "Roles", err: !reflect.DeepEqual(u.Roles, uRes.Roles)},
					{field: "ProfilName", err: u.ProfilName != uRes.ProfilName},
					{field: "ShortId", err: u.ShortId != uRes.ShortId},
				}

				for _, utest := range testUserOutput {
					if utest.err {
						t.Errorf("field output invalid: %s", utest.field)
					}
				}
			}
		})
	}
}

type MeResponse struct {
	Me user.User
}
