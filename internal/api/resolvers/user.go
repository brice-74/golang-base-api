package resolvers

import (
	"context"
	"time"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/domains/user"
	"github.com/brice-74/golang-base-api/internal/utils"
	"github.com/brice-74/golang-base-api/pkg/validator"
	"github.com/graph-gophers/graphql-go"
)

func (r Root) Me(ctx context.Context) (*UserAccountResolver, error) {
	c := r.App.ClientFromContext(ctx)

	if c.User.IsAnonymous() {
		return nil, resolverErrUnauthorized(nil)
	}

	return &UserAccountResolver{app: r.App, user: *c.User}, nil
}

type UserAccountResolver struct {
	app  *application.Application
	user user.User
}

func (r UserAccountResolver) ID() graphql.ID {
	return graphql.ID(r.user.ID)
}

func (r UserAccountResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.user.CreatedAt}
}

func (r UserAccountResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: r.user.UpdatedAt}
}

func (r UserAccountResolver) Active() bool {
	return r.user.DeactivatedAt.IsZero()
}

func (r UserAccountResolver) Email() string {
	return r.user.Email
}

func (r UserAccountResolver) Password() string {
	return r.user.Password
}

func (r UserAccountResolver) Roles() user.Roles {
	return r.user.Roles
}

func (r UserAccountResolver) ProfilName() string {
	return r.user.ProfilName
}

func (r UserAccountResolver) ShortId() string {
	return r.user.ShortId
}

func (r Root) SessionsFromAuth(
	ctx context.Context,
	params SessionListParams,
) (*SessionListResolver, error) {
	c := r.App.ClientFromContext(ctx)

	if c.User.IsAnonymous() {
		return nil, resolverErrUnauthorized(nil)
	}

	v := validator.New()

	qp := utils.QueryParams{
		Sort: params.Sort,
		SortableFields: []string{
			"deactivatedAt",
			"-deactivatedAt",
		},
		Offset: int(params.Offset),
		Limit:  int(params.Limit),
	}

	if qp.Validate(v); !v.Valid() {
		return nil, validatorError{Errors: v.Errors}
	}

	sessions, total, err := r.App.Models.User.GetAllSession(
		qp,
		user.GetAllSessionIncludeFilters{
			States:  params.Include.States,
			UserIds: []string{c.User.ID},
		},
	)
	if err != nil {
		return nil, err
	}

	var sr []SessionResolver
	for _, s := range sessions {
		sr = append(sr, SessionResolver{app: r.App, session: *s})
	}

	return &SessionListResolver{total: total, resolvers: sr}, nil
}

type SessionListParams struct {
	ResolverParams
	Include SessionListIncludeFiltersInput
}

type SessionListIncludeFiltersInput struct {
	States []user.SessionActivityState
}

type SessionListResolver struct {
	total     int
	resolvers []SessionResolver
}

func (r SessionListResolver) Total() int32 {
	return int32(r.total)
}

func (r SessionListResolver) List() []SessionResolver {
	return r.resolvers
}

type SessionResolver struct {
	app     *application.Application
	session user.Session
}

func (r SessionResolver) ID() graphql.ID {
	return graphql.ID(r.session.ID)
}

func (r SessionResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.session.CreatedAt}
}

func (r SessionResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: r.session.UpdatedAt}
}

func (r SessionResolver) Active() bool {
	return r.session.DeactivatedAt.After(time.Now())
}

func (r SessionResolver) IP() string {
	return r.session.IP
}

func (r SessionResolver) Agent() string {
	return r.session.Agent
}

func (r SessionResolver) UserID() string {
	return r.session.UserID
}
