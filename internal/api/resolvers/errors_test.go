package resolvers

import (
	"errors"
	"testing"

	"github.com/brice-74/golang-base-api/pkg/validator"
	"github.com/google/go-cmp/cmp"
)

func TestResolverErrNotFound(t *testing.T) {
	expect := resolverError{
		Code:       errNotFound,
		StatusCode: 404,
	}

	t.Run("Custom error message", func(t *testing.T) {
		err := errors.New("custom error message")
		got := resolverErrNotFound(err)
		expect.Message = err.Error()
		if got != expect {
			t.Fatalf("got resolver error: %+v, expect: %+v", got, expect)
		}
	})

	t.Run("Default error message", func(t *testing.T) {
		got := resolverErrNotFound(nil)
		expect.Message = "Ressource could not be found"
		if got != expect {
			t.Fatalf("got resolver error: %+v, expect: %+v", got, expect)
		}
	})
}

func TestResolverErrUnauthorized(t *testing.T) {
	expect := resolverError{
		Code:       errUnauthorized,
		StatusCode: 401,
	}

	t.Run("Custom error message", func(t *testing.T) {
		err := errors.New("custom error message")
		got := resolverErrUnauthorized(err)
		expect.Message = err.Error()
		if got != expect {
			t.Fatalf("got resolver error: %+v, expect: %+v", got, expect)
		}
	})

	t.Run("Default error message", func(t *testing.T) {
		got := resolverErrUnauthorized(nil)
		expect.Message = "Unauthorized access"
		if got != expect {
			t.Fatalf("got resolver error: %+v, expect: %+v", got, expect)
		}
	})
}

func TestResolverErrDatabaseOperation(t *testing.T) {
	expect := resolverError{
		Code:       errDatabaseOperation,
		StatusCode: 500,
	}

	t.Run("Custom error message", func(t *testing.T) {
		err := errors.New("custom error message")
		got := resolverErrDatabaseOperation(err)
		expect.Message = err.Error()
		if got != expect {
			t.Fatalf("got resolver error: %+v, expect: %+v", got, expect)
		}
	})

	t.Run("Default error message", func(t *testing.T) {
		got := resolverErrDatabaseOperation(nil)
		expect.Message = "Database operation error"
		if got != expect {
			t.Fatalf("got resolver error: %+v, expect: %+v", got, expect)
		}
	})
}

func TestResolverError(t *testing.T) {
	e := resolverError{
		Code:    "err",
		Message: "message",
	}

	expect := "error [err]: message"

	if e.Error() != expect {
		t.Fatalf("got string error: %s, expect: %s", e.Error(), expect)
	}
}

func TestResolverExtensions(t *testing.T) {
	e := resolverError{
		Code:       "err",
		StatusCode: 500,
		Message:    "message",
	}

	expect := map[string]interface{}{
		"statusCode": 500,
		"code":       "err",
		"message":    "message",
	}

	if diff := cmp.Diff(expect, e.Extensions()); diff != "" {
		t.Fatal(diff)
	}
}

func TestValidatorError(t *testing.T) {
	e := validatorError{}

	expect := "validation error [ValidatorError]"

	if e.Error() != expect {
		t.Fatalf("got string error: %s, expect: %s", e.Error(), expect)
	}
}

func TestValidatorExtensions(t *testing.T) {
	errs := validator.Errors{
		"key": []string{
			"error",
		},
	}

	e := validatorError{
		Errors: errs,
	}

	expect := map[string]interface{}{
		"statusCode": 422,
		"code":       "ValidatorError",
		"errors":     errs,
	}

	if diff := cmp.Diff(expect, e.Extensions()); diff != "" {
		t.Fatal(diff)
	}
}
