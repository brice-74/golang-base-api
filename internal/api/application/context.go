package application

import (
	"context"

	"github.com/brice-74/golang-base-api/internal/domains/user"
)

type contextKey string

const (
	ClientCtxKey = contextKey("client")
)

// ContextWithClient returns a new ClientCtx instance added in the context.
func (app *Application) ContextWithClient(ctx context.Context, user *ClientCtx) context.Context {
	return context.WithValue(ctx, ClientCtxKey, user)
}

// ClientFromContext retrieves the ClientCtx struct from the request context.
func (app *Application) ClientFromContext(ctx context.Context) *ClientCtx {
	u, ok := ctx.Value(ClientCtxKey).(*ClientCtx)
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
