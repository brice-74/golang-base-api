package application

import (
	"context"
	"testing"

	"github.com/brice-74/golang-base-api/internal/domains/user"
	"github.com/google/go-cmp/cmp"
)

func TestContextWithUser(t *testing.T) {
	app := Application{}

	u := &UserCtx{
		User: &user.User{
			ID: "1234",
		},
		Client: &Client{
			SessionID: "5678",
			IP:        "0.0.0.0",
			Agent:     "agent",
		},
	}

	ctx := app.ContextWithUser(context.Background(), u)

	got := ctx.Value(userCtxKey)

	if diff := cmp.Diff(got, u); diff != "" {
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

		app.UserFromContext(context.Background())
	})

	t.Run("should return user context", func(t *testing.T) {
		u := &UserCtx{
			User: &user.User{
				ID: "1234",
			},
			Client: &Client{
				SessionID: "5678",
				IP:        "0.0.0.0",
				Agent:     "agent",
			},
		}

		ctx := context.WithValue(context.Background(), userCtxKey, u)

		got := app.UserFromContext(ctx)

		if diff := cmp.Diff(got, u); diff != "" {
			t.Fatal(diff)
		}
	})
}
