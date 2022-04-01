package resolvers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/api/resolvers"
	"github.com/brice-74/golang-base-api/internal/domains/user"
	"github.com/brice-74/golang-base-api/internal/testutils"
	"github.com/brice-74/golang-base-api/internal/testutils/factory"
	"github.com/brice-74/golang-base-api/internal/utils"
	"github.com/brice-74/golang-base-api/pkg/validator"
	"github.com/google/go-cmp/cmp"
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
					{field: "DeactivatedAt", err: !uRes.Active},
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
	Me UserResponse
}

func TestSessionsFromAuth(t *testing.T) {
	var (
		db     = testutils.PrepareDB(t)
		app    = testutils.NewApplication(db)
		schema = testutils.ParseTestSchema(app)
		fac    = factory.New(t, db)
	)

	u := fac.CreateUserAccount(nil)
	sActiv := fac.CreateUserSession(&user.Session{
		DeactivatedAt: time.Now().Add(time.Hour).UTC().Round(time.Second),
		UserID:        u.ID,
	})
	sExp := fac.CreateUserSession(&user.Session{
		DeactivatedAt: time.Date(2021, 0, 0, 0, 0, 0, 0, time.UTC),
		UserID:        u.ID,
	})

	var sessionToSessionResponse = func(s user.Session) SessionResponse {
		return SessionResponse{
			ID:        s.ID,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
			Active:    s.DeactivatedAt.After(time.Now()),
			IP:        s.IP,
			Agent:     s.Agent,
			UserID:    s.UserID,
		}
	}

	sActivRes := sessionToSessionResponse(*sActiv)
	sExpRes := sessionToSessionResponse(*sExp)

	params := utils.QueryParams{
		Offset: 0,
		Limit:  20,
		Sort:   "deactivatedAt",
		SortableFields: []string{
			"deactivatedAt",
		},
	}

	paramsErr := utils.QueryParams{}

	var queryString = func(params utils.QueryParams, inc resolvers.SessionListIncludeFiltersInput) string {
		return fmt.Sprintf(`
			{
				sessionsFromAuth(
					offset: %d,
					limit: %d,
					sort: "%s",
					include: {
						states: %s
					},
				) {
					total
					list {
						id
						createdAt
						updatedAt
						active
						ip
						agent
						userId
					}
				}
			}
		`,
			params.Offset, params.Limit, params.Sort,
			strings.ReplaceAll(fmt.Sprintf("%s", inc.States), " ", ", "),
		)
	}

	tests := []struct {
		title          string
		gqltest        *gqltesting.Test
		expectSessions []*SessionResponse
		expectError    *testutils.ExpectResolverError
	}{
		{
			title: "Should return all sessions from user",
			gqltest: &gqltesting.Test{
				Context: app.ContextWithClient(context.Background(), &application.ClientCtx{
					User: u,
				}),
				Schema: schema,
				Query:  queryString(params, resolvers.SessionListIncludeFiltersInput{}),
			},
			expectSessions: []*SessionResponse{&sExpRes, &sActivRes},
		},
		{
			title: "Should return all active sessions from user",
			gqltest: &gqltesting.Test{
				Context: app.ContextWithClient(context.Background(), &application.ClientCtx{
					User: u,
				}),
				Schema: schema,
				Query: queryString(params, resolvers.SessionListIncludeFiltersInput{
					States: []user.SessionActivityState{user.SessionActive},
				}),
			},
			expectSessions: []*SessionResponse{&sActivRes},
		},
		{
			title: "Should be unauthorize",
			gqltest: &gqltesting.Test{
				Context: app.ContextWithClient(context.Background(), &application.ClientCtx{
					User: user.AnonymousUser,
				}),
				Schema: schema,
				Query:  queryString(params, resolvers.SessionListIncludeFiltersInput{}),
			},
			expectError: &testutils.ExpectResolverError{
				Msg: "error [Unauthorized]: Unauthorized access",
				Extensions: map[string]interface{}{
					"code":       "Unauthorized",
					"statusCode": 401,
					"message":    "Unauthorized access",
				},
			},
		},
		{
			title: "Should be query params error",
			gqltest: &gqltesting.Test{
				Context: app.ContextWithClient(context.Background(), &application.ClientCtx{
					User: u,
				}),
				Schema: schema,
				Query:  queryString(paramsErr, resolvers.SessionListIncludeFiltersInput{}),
			},
			expectError: &testutils.ExpectResolverError{
				Msg: "validation error [ValidatorError]",
				Extensions: map[string]interface{}{
					"code":       "ValidatorError",
					"statusCode": 422,
					"errors": validator.Errors{
						"limit": []string{"must be greater than zero"},
						"sort":  []string{"invalid sort value"},
					},
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
				var res SessionsFromAuthResponse

				data, _ := result.Data.MarshalJSON()
				if err := json.Unmarshal(data, &res); err != nil {
					t.Fatal(err)
				}

				lens := len(tt.expectSessions)

				if lens != res.SessionsFromAuth.Total {
					t.Errorf("expect total: %d, got: %d", lens, res.SessionsFromAuth.Total)
				}

				if diff := cmp.Diff(tt.expectSessions, res.SessionsFromAuth.List); diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}

type SessionsFromAuthResponse struct {
	SessionsFromAuth SessionListResponse
}

type SessionListResponse struct {
	Total int
	List  []*SessionResponse
}

type SessionResponse struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Active    bool
	IP        string
	Agent     string
	UserID    string
}

type UserResponse struct {
	ID         string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Active     bool
	Email      string
	Password   string
	Roles      user.Roles
	ProfilName string
	ShortId    string
}
