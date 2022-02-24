package application

import (
	"context"
	"testing"

	"github.com/brice-74/golang-base-api/internal/domains/user"
	"github.com/google/go-cmp/cmp"
)

func TestContextWithUser(t *testing.T) {
	app := Application{}

	c := &ClientCtx{
		User: &user.User{
			ID: "1234",
		},
		Agent: &Agent{
			IP:    "0.0.0.0",
			Agent: "agent",
		},
		Session: &user.Session{
			ID: "5678",
		},
	}

	ctx := app.ContextWithClient(context.Background(), c)

	got := ctx.Value(clientCtxKey)

	if diff := cmp.Diff(got, c); diff != "" {
		t.Fatal(diff)
	}
}

func TestUserFromContext(t *testing.T) {
	app := Application{}

	t.Run("should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("got no panic, expected panic")
			}
		}()

		app.ClientFromContext(context.Background())
	})

	t.Run("should return user context", func(t *testing.T) {
		c := &ClientCtx{
			User: &user.User{
				ID: "1234",
			},
			Agent: &Agent{
				IP:    "0.0.0.0",
				Agent: "agent",
			},
			Session: &user.Session{
				ID: "5678",
			},
		}

		ctx := context.WithValue(context.Background(), clientCtxKey, c)

		got := app.ClientFromContext(ctx)

		if diff := cmp.Diff(got, c); diff != "" {
			t.Fatal(diff)
		}
	})
}
