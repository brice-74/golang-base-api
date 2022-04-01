package resolvers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/domains/user"
	"github.com/brice-74/golang-base-api/internal/testutils"
	"github.com/brice-74/golang-base-api/internal/testutils/factory"
	"github.com/brice-74/golang-base-api/internal/testutils/mocks"
	"github.com/brice-74/golang-base-api/pkg/validator"
	"github.com/dgrijalva/jwt-go"
	"github.com/graph-gophers/graphql-go/gqltesting"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
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

	var queryString = func(email, pass, sessionID string) string {
		return fmt.Sprintf(`
			mutation {
				registerUserAccount(input: {
					email: "%s",
					password: "%s",
					profilName: "%s"
				}) {
					id
					createdAt
					updatedAt
					active
					email
					password
					roles
					profilName
					shortId
				}
			}`, email, pass, sessionID,
		)
	}

	tests := []struct {
		title       string
		gqltest     *gqltesting.Test
		expectError *testutils.ExpectResolverError
	}{
		{
			title: "Should insert and return available user",
			gqltest: &gqltesting.Test{
				Schema: schema,
				Query:  queryString(email, password, profilName),
			},
		},
		{
			title: "Should return database error conflict",
			gqltest: &gqltesting.Test{
				Schema: schema,
				Query:  queryString(email, password, profilName),
			},
			expectError: &testutils.ExpectResolverError{
				Msg: "error [DatabaseOperationError]: Duplicate email",
				Extensions: map[string]interface{}{
					"code":       "DatabaseOperationError",
					"statusCode": 500,
					"message":    "Duplicate email",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			result := tt.gqltest.Schema.Exec(context.Background(), tt.gqltest.Query, "", nil)

			if tt.expectError != nil {
				testutils.TestGqlError(t, result.Errors[0], tt.expectError)
			} else {
				var res RegisterUserAccountResponse

				data, _ := result.Data.MarshalJSON()
				if err := json.Unmarshal(data, &res); err != nil {
					t.Fatal(err)
				}

				var u = res.RegisterUserAccount

				if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
					t.Fatal("incorrect password from resolver")
				}

				testUserOutput := []struct {
					field string
					err   bool
				}{
					{field: "ID", err: u.ID == ""},
					{field: "CreatedAt", err: u.CreatedAt.IsZero()},
					{field: "UpdatedAt", err: u.UpdatedAt.IsZero()},
					{field: "Active", err: !u.Active},
					{field: "Email", err: u.Email != email},
					{field: "Roles", err: u.Roles[0] != user.RoleUser},
					{field: "ProfilName", err: u.ProfilName == ""},
					{field: "ShortId", err: u.ShortId == ""},
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

type RegisterUserAccountResponse struct {
	RegisterUserAccount UserResponse
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
		expectError *testutils.ExpectResolverError
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
			expectError: &testutils.ExpectResolverError{
				Msg: "validation error [ValidatorError]",
				Extensions: map[string]interface{}{
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
			expectError: &testutils.ExpectResolverError{
				Msg: "error [NotFoundError]: User not found",
				Extensions: map[string]interface{}{
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
			expectError: &testutils.ExpectResolverError{
				Msg: "error [Unauthorized]: incorrect password",
				Extensions: map[string]interface{}{
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
				testutils.TestGqlError(t, result.Errors[0], tt.expectError)
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

func TestRefreshUserAccount(t *testing.T) {
	var (
		db     = testutils.PrepareDB(t)
		app    = testutils.NewApplication(db)
		schema = testutils.ParseTestSchema(app)
		fac    = factory.New(t, db)
	)

	var queryString = func(token string) string {
		return fmt.Sprintf(`
			mutation {
				refreshUserAccount(
					token: "%s"
				) {
					access
					refresh
				}
			}`, token,
		)
	}

	var queryContext = app.ContextWithClient(context.Background(), &application.ClientCtx{
		Agent: &application.Agent{
			IP:    "0.0.0.0",
			Agent: "agent",
		},
	})

	u := fac.CreateUserAccount(&user.User{
		Email:    "test@test.com",
		Password: "Test123!",
	})

	s := fac.CreateUserSession(&user.Session{
		UserID: u.ID,
	})

	exp := time.Now().Add(time.Minute * 3)

	uuid := uuid.NewV4().String()
	goodClaims := mocks.CreateClaims(u.ID, s.ID, exp)
	badUuidClaims := mocks.CreateClaims("", "", exp)
	badIdsClaims := mocks.CreateClaims(uuid, uuid, exp)

	goodToken := mocks.CreateToken(t, jwt.SigningMethodHS256, goodClaims, app.Config.JWT.Refresh.Secret)
	badClaimsToken := mocks.CreateToken(t, jwt.SigningMethodHS256, nil, app.Config.JWT.Refresh.Secret)
	badUuidToken := mocks.CreateToken(t, jwt.SigningMethodHS256, badUuidClaims, app.Config.JWT.Refresh.Secret)
	badIdsToken := mocks.CreateToken(t, jwt.SigningMethodHS256, badIdsClaims, app.Config.JWT.Refresh.Secret)

	tests := []struct {
		title       string
		gqltest     *gqltesting.Test
		expectError *testutils.ExpectResolverError
	}{
		{
			title: "Should return available tokens",
			gqltest: &gqltesting.Test{
				Schema:  schema,
				Context: queryContext,
				Query:   queryString(goodToken),
			},
		},
		{
			title: "Should return claims not found from token",
			gqltest: &gqltesting.Test{
				Schema:  schema,
				Context: queryContext,
				Query:   queryString(badClaimsToken),
			},
			expectError: &testutils.ExpectResolverError{
				Msg: "Required claims from token not found",
			},
		},
		{
			title: "Should return server error",
			gqltest: &gqltesting.Test{
				Schema:  schema,
				Context: queryContext,
				Query:   queryString(badUuidToken),
			},
			expectError: &testutils.ExpectResolverError{
				Msg: `error [DatabaseOperationError]: pq: invalid input syntax for type uuid: ""`,
				Extensions: map[string]interface{}{
					"code":       "DatabaseOperationError",
					"statusCode": 500,
					"message":    `pq: invalid input syntax for type uuid: ""`,
				},
			},
		},
		{
			title: "Should return not found",
			gqltest: &gqltesting.Test{
				Schema:  schema,
				Context: queryContext,
				Query:   queryString(badIdsToken),
			},
			expectError: &testutils.ExpectResolverError{
				Msg: "error [NotFoundError]: User or user session not found",
				Extensions: map[string]interface{}{
					"code":       "NotFoundError",
					"statusCode": 404,
					"message":    "User or user session not found",
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
				var res RefreshUserAccountResponse

				data, _ := result.Data.MarshalJSON()
				if err := json.Unmarshal(data, &res); err != nil {
					t.Fatal(err)
				}

				if _, err := application.VerifyToken(res.RefreshUserAccount.Access, app.Config.JWT.Access.Secret); err != nil {
					t.Fatalf("Access token verification fail: %s", err.Error())
				}

				if _, err := application.VerifyToken(res.RefreshUserAccount.Refresh, app.Config.JWT.Refresh.Secret); err != nil {
					t.Fatalf("Refresh token verification fail: %s", err.Error())
				}
			}
		})
	}
}

type RefreshUserAccountResponse struct {
	RefreshUserAccount user.Tokens
}
