package user_test

import (
	"testing"
	"time"

	"github.com/brice-74/golang-base-api/internal/domains/user"
	"github.com/brice-74/golang-base-api/internal/testutils"
	"github.com/brice-74/golang-base-api/internal/testutils/factory"
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

	got, err := m.GetById(u.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(u, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestGetByEmail(t *testing.T) {
	var (
		db  = testutils.PrepareDB(t)
		m   = user.Model{DB: db}
		fac = factory.New(t, db)
	)

	u := fac.CreateUserAccount(nil)

	got, err := m.GetByEmail(u.Email)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(u, got); diff != "" {
		t.Fatal(diff)
	}
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
}
