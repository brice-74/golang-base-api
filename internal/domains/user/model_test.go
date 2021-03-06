package user_test

import (
	"testing"
	"time"

	"github.com/brice-74/golang-base-api/internal/domains/user"
	"github.com/brice-74/golang-base-api/internal/testutils"
	"github.com/brice-74/golang-base-api/internal/testutils/factory"
	"github.com/brice-74/golang-base-api/internal/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/twinj/uuid"
)

func TestExistEmail(t *testing.T) {
	var (
		db  = testutils.PrepareDB(t)
		m   = user.Model{DB: db}
		fac = factory.New(t, db)
	)

	u := fac.CreateUserAccount(nil)

	find, err := m.ExistEmail(u.Email)
	if err != nil {
		t.Fatal(err)
	}

	if !find {
		t.Fatalf("user email must be find: %s", u.Email)
	}
}

func TestInsertRegisteredUserAccount(t *testing.T) {
	var (
		db = testutils.PrepareDB(t)
		m  = user.Model{DB: db}
	)

	u := &user.User{
		Email:      "test@test.com",
		Password:   "passtest",
		Roles:      user.Roles{user.RoleUser},
		ProfilName: "profiletest",
		ShortId:    "shortid",
	}

	t.Run("should insert user", func(t *testing.T) {
		if err := m.InsertRegisteredUserAccount(u); err != nil {
			t.Fatalf("got an error during insert user execution: %s", err)
		}

		if u.ID == "" {
			t.Error("got id zero value instead of a generated uuid")
		}

		if u.CreatedAt.IsZero() {
			t.Error("got CreatedAt zero value instead of a real date")
		}

		if u.UpdatedAt.IsZero() {
			t.Error("got UpdatedAt zero value instead of a real date")
		}
	})

	t.Run("should return duplicate user error", func(t *testing.T) {
		err := m.InsertRegisteredUserAccount(u)
		if err == nil {
			t.Fatal("got nil error, expect available error")
		}

		if err.Error() != user.ErrDuplicateEmail.Error() {
			t.Fatalf("got: %s, expect: %s", err.Error(), user.ErrDuplicateEmail.Error())
		}
	})
}

func TestInsertOrUpdateUserSession(t *testing.T) {
	var (
		db  = testutils.PrepareDB(t)
		m   = user.Model{DB: db}
		fac = factory.New(t, db)
	)

	u := fac.CreateUserAccount(nil)

	s := &user.Session{
		ID:            uuid.NewV4().String(),
		DeactivatedAt: time.Now(),
		IP:            "0.0.0.0",
		Agent:         "agent",
		UserID:        u.ID,
	}

	if err := m.InsertOrUpdateUserSession(s); err != nil {
		t.Fatalf("got an error during insert user session execution: %s", err)
	}

	if s.CreatedAt.IsZero() {
		t.Error("got CreatedAt zero value instead of a real date")
	}

	if s.UpdatedAt.IsZero() {
		t.Error("got UpdatedAt zero value instead of a real date")
	}
}

func TestGetById(t *testing.T) {
	var (
		db  = testutils.PrepareDB(t)
		m   = user.Model{DB: db}
		fac = factory.New(t, db)
	)

	u := fac.CreateUserAccount(nil)

	t.Run("should find user", func(t *testing.T) {
		got, err := m.GetById(u.ID)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(u, got); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("shouldn't find user", func(t *testing.T) {
		_, err := m.GetById(uuid.NewV4().String())

		if err == nil {
			t.Fatal("got nil error, expect available error")
		}

		if err.Error() != user.ErrNotFoundUser.Error() {
			t.Fatalf("got: %s, expect: %s", err.Error(), user.ErrNotFoundUser.Error())
		}
	})
}

func TestGetByEmail(t *testing.T) {
	var (
		db  = testutils.PrepareDB(t)
		m   = user.Model{DB: db}
		fac = factory.New(t, db)
	)

	u := fac.CreateUserAccount(nil)

	t.Run("should find user", func(t *testing.T) {
		got, err := m.GetByEmail(u.Email)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(u, got); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("shouldn't find user", func(t *testing.T) {
		_, err := m.GetByEmail("test@test.com")

		if err == nil {
			t.Fatal("got nil error, expect available error")
		}

		if err.Error() != user.ErrNotFoundUser.Error() {
			t.Fatalf("got: %s, expect: %s", err.Error(), user.ErrNotFoundUser.Error())
		}
	})
}

func TestGetSessionByID(t *testing.T) {
	var (
		db  = testutils.PrepareDB(t)
		m   = user.Model{DB: db}
		fac = factory.New(t, db)
	)

	s := fac.CreateUserSession(nil)

	got, err := m.GetSessionByID(s.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(s, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestGetUserAndSession(t *testing.T) {
	var (
		db  = testutils.PrepareDB(t)
		m   = user.Model{DB: db}
		fac = factory.New(t, db)
	)

	u := fac.CreateUserAccount(nil)
	s := fac.CreateUserSession(&user.Session{
		UserID: u.ID,
	})

	t.Run("should find user and session", func(t *testing.T) {
		gotu, gots, err := m.GetUserAndSession(u.ID, s.ID)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(u, gotu); diff != "" {
			t.Fatal(diff)
		}

		if diff := cmp.Diff(s, gots); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("shouldn't find user and session", func(t *testing.T) {
		_, _, err := m.GetUserAndSession(u.ID, uuid.NewV4().String())

		if err == nil {
			t.Fatal("got nil error, expect available error")
		}

		if err.Error() != user.ErrNotFoundUserAndSession.Error() {
			t.Fatalf("got: %s, expect: %s", err.Error(), user.ErrNotFoundUserAndSession.Error())
		}
	})
}

func TestGetAllSession(t *testing.T) {
	var (
		db  = testutils.PrepareDB(t)
		m   = user.Model{DB: db}
		fac = factory.New(t, db)
	)

	sActiv := fac.CreateUserSession(&user.Session{
		DeactivatedAt: time.Now().Add(time.Hour).UTC().Round(time.Second),
	})
	sExp := fac.CreateUserSession(&user.Session{
		DeactivatedAt: time.Date(2021, 0, 0, 0, 0, 0, 0, time.UTC),
	})

	params := utils.QueryParams{
		Offset: 0,
		Limit:  20,
		Sort:   "deactivatedAt",
		SortableFields: []string{
			"deactivatedAt",
		},
	}

	tests := []struct {
		title          string
		include        user.GetAllSessionIncludeFilters
		expectSessions []*user.Session
	}{
		{
			title:          "Should return all sessions",
			include:        user.GetAllSessionIncludeFilters{},
			expectSessions: []*user.Session{sExp, sActiv},
		},
		{
			title:          "Should return sessions by user id",
			include:        user.GetAllSessionIncludeFilters{UserIds: []string{sExp.UserID}},
			expectSessions: []*user.Session{sExp},
		},
		{
			title:          "Should return active sessions",
			include:        user.GetAllSessionIncludeFilters{States: []user.SessionActivityState{user.SessionActive}},
			expectSessions: []*user.Session{sActiv},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			sessions, total, err := m.GetAllSession(params, tt.include)
			if err != nil {
				t.Fatal(err)
			}

			lens := len(tt.expectSessions)

			if total != lens {
				t.Errorf("expect total: %d, got: %d", lens, total)
			}

			if diff := cmp.Diff(tt.expectSessions, sessions); diff != "" {
				t.Error(diff)
			}
		})
	}

}
