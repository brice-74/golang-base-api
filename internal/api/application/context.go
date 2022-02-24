package application

import (
	"context"

	"github.com/brice-74/golang-base-api/internal/domains/user"
)

type contextKey string

const (
	clientCtxKey = contextKey("client")
)

// ContextWithUser returns a new User instance added in the context.
func (app *Application) ContextWithClient(ctx context.Context, user *ClientCtx) context.Context {
	return context.WithValue(ctx, clientCtxKey, user)
}

// UserFromContext retrieves the User struct from the request context.
func (app *Application) ClientFromContext(ctx context.Context) *ClientCtx {
	u, ok := ctx.Value(clientCtxKey).(*ClientCtx)
	if !ok {
		panic("missing user value in request context")
	}

	return u
}

type ClientCtx struct {
	Agent   *Agent
	User    *user.User
	Session *user.Session
}

type Agent struct {
	IP    string
	Agent string
}
