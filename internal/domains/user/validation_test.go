package user_test

import (
	"testing"

	"github.com/brice-74/golang-base-api/internal/domains/user"
	"github.com/brice-74/golang-base-api/pkg/validator"
	"github.com/google/go-cmp/cmp"
	"github.com/jaswdr/faker"
)

func TestValidateProfilNameEntry(t *testing.T) {
	var (
		v      = validator.New()
		u      = user.User{}
		errKey = "profil name"
	)

	t.Run("should be provided", func(t *testing.T) {
		u.ProfilName = ""
		u.ValidateProfilNameEntry(v)

		got := v.Errors[errKey][0]
		expect := "must be provided"

		if got != expect {
			t.Fatalf("unexpected error got: %s, expected: %s", got, expect)
		}
	})

	v.Errors = make(validator.Errors)
	t.Run("max 32 characters", func(t *testing.T) {
		for {
			u.ProfilName = faker.New().Lorem().Text(50)
			if len(u.ProfilName) > 32 {
				break
			}
		}

		u.ValidateProfilNameEntry(v)

		got := v.Errors[errKey][0]
		expect := "must have maximum of 32 characters"

		if got != expect {
			t.Fatalf("unexpected error got: %s, expected: %s", got, expect)
		}
	})

	v.Errors = make(validator.Errors)
	t.Run("should be ok", func(t *testing.T) {
		u.ProfilName = "abcde"
		u.ValidateProfilNameEntry(v)

		if !v.Valid() {
			t.Fatal("validator should be valid")
		}
	})
}

func TestValidatePasswordEntry(t *testing.T) {
	var (
		v      = validator.New()
		u      = user.User{}
		errKey = "password"
	)

	t.Run("min 8 characters, 1 lower, upper, number, special", func(t *testing.T) {
		u.Password = ""
		u.ValidatePasswordEntry(v)

		got := v.Errors[errKey]

		expect := make(validator.Errors)[errKey]
		expect = []string{
			"must be provided",
			"must have minimum of 8 characters",
			"must have minimum of 1 lowercase",
			"must have minimum of 1 uppercase",
			"must have minimum of 1 number",
			"must have minimum of 1 special character",
		}

		if diff := cmp.Diff(got, expect); diff != "" {
			t.Fatal(diff)
		}
	})

	v.Errors = make(validator.Errors)
	t.Run("max 255 characters", func(t *testing.T) {
		for {
			u.Password = faker.New().Lorem().Text(300) + "A1!"
			if len(u.Password) > 255 {
				break
			}
		}

		u.ValidatePasswordEntry(v)

		got := v.Errors[errKey][0]
		expect := "must have maximum of 255 characters"

		if got != expect {
			t.Fatalf("unexpected error got: %s, expected: %s", got, expect)
		}
	})

	v.Errors = make(validator.Errors)
	t.Run("should be ok", func(t *testing.T) {
		u.Password = "Test123!"
		u.ValidatePasswordEntry(v)

		if !v.Valid() {
			t.Fatal("validator should be valid")
		}
	})
}

func TestValidateEmailEntry(t *testing.T) {
	var (
		v      = validator.New()
		u      = user.User{}
		errKey = "email"
	)

	t.Run("should be provided and valid", func(t *testing.T) {
		u.Email = ""
		u.ValidateEmailEntry(v)

		got := v.Errors[errKey]

		expect := make(validator.Errors)[errKey]
		expect = []string{
			"must be provided",
			"must be a valid address",
		}

		if diff := cmp.Diff(got, expect); diff != "" {
			t.Fatal(diff)
		}
	})

	v.Errors = make(validator.Errors)
	t.Run("max 255 characters", func(t *testing.T) {
		for {
			u.Email = faker.New().Lorem().Text(300) + "@test.com"
			if len(u.Email) > 255 {
				break
			}
		}

		u.ValidateEmailEntry(v)

		got := v.Errors[errKey][0]
		expect := "must have maximum of 255 characters"

		if got != expect {
			t.Fatalf("unexpected error got: %s, expected: %s", got, expect)
		}
	})

	v.Errors = make(validator.Errors)
	t.Run("should be ok", func(t *testing.T) {
		u.Email = "test@test.com"
		u.ValidateEmailEntry(v)

		if !v.Valid() {
			t.Fatal("validator should be valid")
		}
	})
}
