package user_test

import (
	"testing"
	"time"

	"github.com/brice-74/golang-base-api/internal/domains/user"
)

func TestIsActiveSession(t *testing.T) {
	var s = user.Session{
		DeactivatedAt: time.Date(2021, 1, 1, 1, 1, 1, 1, time.UTC),
	}

	if s.IsActive() {
		t.Fatalf("Session with deactivated date: %s must be inactive", s.DeactivatedAt.Format("2006-01-02"))
	}
}
