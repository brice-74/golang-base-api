package application

import (
	"context"

	"github.com/brice-74/golang-base-api/internal/domains/user"
)

type contextKey string

const (
	userCtxKey = contextKey("user")
)

// ContextWithUser returns a new User instance added in the context.
func (app *Application) ContextWithUser(ctx context.Context, user *UserCtx) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

// UserFromContext retrieves the User struct from the request context.
func (app *Application) UserFromContext(ctx context.Context) *UserCtx {
	u, ok := ctx.Value(userCtxKey).(*UserCtx)
	if !ok {
		panic("missing user value in request context")
	}

	return u
}

type UserCtx struct {
	User      *user.User
	Token     *TokenCtx
	SessionID string
}

type TokenCtx struct {
	IsAccess bool
}
