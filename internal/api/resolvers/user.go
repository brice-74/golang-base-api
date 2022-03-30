package resolvers

import (
	"context"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/domains/user"
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
