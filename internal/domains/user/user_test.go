package user_test

import (
	"testing"

	"github.com/brice-74/golang-base-api/internal/domains/user"
)

func TestIsAnonymous(t *testing.T) {
	u := user.AnonymousUser

	if !u.IsAnonymous() {
		t.Fatalf("User should be anonymous: %+v", u.Roles)
	}
}
