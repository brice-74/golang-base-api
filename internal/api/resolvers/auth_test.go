package resolvers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/domains/user"
	"github.com/brice-74/golang-base-api/internal/testutils"
	"github.com/brice-74/golang-base-api/internal/testutils/factory"
	"github.com/brice-74/golang-base-api/pkg/validator"
	"github.com/google/go-cmp/cmp"
	"github.com/graph-gophers/graphql-go/gqltesting"
	"github.com/twinj/uuid"
)

func TestRegisterUserAccount(t *testing.T) {
	var (
		db     = testutils.PrepareDB(t)
		app    = testutils.NewApplication(db)
		schema = testutils.ParseTestSchema(app)
	)

	const (
		email      string = "test@test.com"
		password   string = "passWORD123!"
		profilName string = "name"
	)

	test := &gqltesting.Test{
		Schema: schema,
		Query: fmt.Sprintf(`
			mutation {
				registerUserAccount(input: {
					email: "%s",
					password: "%s",
					profilName: "%s"
				}) {
					email
					profilName
					roles
				}
			}`, email, password, profilName),
		ExpectedResult: fmt.Sprintf(`
			{
				"registerUserAccount": {
					"email": "%s",
					"profilName": "%s",
					"roles": [
						"ROLE_USER"
					]
				}
			}`, email, profilName),
	}

	gqltesting.RunTest(t, test)
}

func TestLoginUserAccount(t *testing.T) {
	var (
		db        = testutils.PrepareDB(t)
		app       = testutils.NewApplication(db)
		schema    = testutils.ParseTestSchema(app)
		fac       = factory.New(t, db)
		sessionID = uuid.NewV4().String()
	)

	var queryString = func(email, pass, sessionID string) string {
		return fmt.Sprintf(`
			mutation {
				loginUserAccount(
					email: "%s",
					password: "%s",
					sessionID: "%s"
				) {
					access
					refresh
				}
			}`, email, pass, sessionID,
		)
	}

	var queryContext = app.ContextWithClient(context.Background(), &application.ClientCtx{
		Agent: &application.Agent{
			IP:    "0.0.0.0",
			Agent: "agent",
		},
	})

	const strPass = "Test123!"
	u := fac.CreateUserAccount(&user.User{
		Email:    "test@test.com",
		Password: strPass,
	})

	tests := []struct {
		title       string
		gqltest     *gqltesting.Test
		expectError *ExpectResolverError
	}{
		{
			title: "Should return available tokens",
			gqltest: &gqltesting.Test{
				Schema:  schema,
				Context: queryContext,
				Query:   queryString(u.Email, strPass, sessionID),
			},
		},
		{
			title: "Should return validatorError",
			gqltest: &gqltesting.Test{
				Schema:  schema,
				Context: queryContext,
				Query:   queryString("bad email", "bad pass", sessionID),
			},
			expectError: &ExpectResolverError{
				msg: "validation error [ValidatorError]",
				extensions: map[string]interface{}{
					"code":       "ValidatorError",
					"statusCode": 422,
					"errors": validator.Errors{
						"email": []string{"must be a valid address"},
						"password": []string{
							"must have minimum of 1 uppercase",
							"must have minimum of 1 number",
							"must have minimum of 1 special character",
						},
					},
				},
			},
		},
		{
			title: "Should return not found user",
			gqltest: &gqltesting.Test{
				Schema:  schema,
				Context: queryContext,
				Query:   queryString("unknow@email.com", strPass, sessionID),
			},
			expectError: &ExpectResolverError{
				msg: "error [NotFoundError]: User not found",
				extensions: map[string]interface{}{
					"code":       "NotFoundError",
					"statusCode": 404,
					"message":    "User not found",
				},
			},
		},
		{
			title: "Should return incorrect password",
			gqltest: &gqltesting.Test{
				Schema:  schema,
				Context: queryContext,
				Query:   queryString(u.Email, "IncorrectPass123!", sessionID),
			},
			expectError: &ExpectResolverError{
				msg: "error [Unauthorized]: incorrect password",
				extensions: map[string]interface{}{
					"code":       "Unauthorized",
					"statusCode": 401,
					"message":    "incorrect password",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			result := tt.gqltest.Schema.Exec(tt.gqltest.Context, tt.gqltest.Query, "", nil)

			if tt.expectError != nil {
				qerr := result.Errors[0]
				if qerr != nil {
					if qerr.Message != tt.expectError.msg {
						t.Errorf("expect query error message: %s, got: %s", tt.expectError.msg, qerr.Message)
					}
					if tt.expectError.extensions != nil {
						if diff := cmp.Diff(tt.expectError.extensions, qerr.Extensions); diff != "" {
							t.Error(diff)
						}
					}
				} else {
					t.Fatal("resolver error expected but not found")
				}
			} else {
				var res LoginUserAccountResponse

				data, _ := result.Data.MarshalJSON()
				if err := json.Unmarshal(data, &res); err != nil {
					t.Fatal(err)
				}

				if _, err := application.VerifyToken(res.LoginUserAccount.Access, app.Config.JWT.Access.Secret); err != nil {
					t.Fatalf("Access token verification fail: %s", err.Error())
				}

				if _, err := application.VerifyToken(res.LoginUserAccount.Refresh, app.Config.JWT.Refresh.Secret); err != nil {
					t.Fatalf("Refresh token verification fail: %s", err.Error())
				}

				s, err := app.Models.User.GetSessionByID(sessionID)
				if err != nil {
					t.Fatalf("error during database session recovery: %s", err.Error())
				}

				if s == nil {
					t.Fatal("expect existing session, got nil")
				}
			}
		})
	}
}

type LoginUserAccountResponse struct {
	LoginUserAccount user.Tokens
}

type ExpectResolverError struct {
	msg        string
	extensions map[string]interface{}
}
