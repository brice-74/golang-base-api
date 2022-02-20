package application

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brice-74/golang-base-api/internal/testutils/mocks"
	"github.com/brice-74/golang-base-api/internal/testutils/require"
)

func TestErrorResponse(t *testing.T) {
	app := Application{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	app.ErrorResponse(w, r, http.StatusNotFound, "not found")

	got := w.Body.String()
	expected := `{"error":"not found"}`

	require.JSONEqual(t, got, expected)
}

func TestServerErrorResponse(t *testing.T) {
	l := mocks.NewLogger()

	app := Application{
		Logger: l,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	app.ServerErrorResponse(w, r, errors.New("server error"))

	got := w.Body.String()
	expected := `{"error":"the server encountered a problem and could not process your request"}`

	require.JSONEqual(t, got, expected)

	if got, expected := w.Code, http.StatusInternalServerError; got != expected {
		t.Fatalf("got status code %d, expected %d", got, expected)
	}

	if !l.PrintErrorCalled {
		t.Error("print error should be called")
	}
}

func TestNotFoundResponse(t *testing.T) {
	app := Application{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	app.NotFoundResponse(w, r)

	got := w.Body.String()
	expected := `{"error":"the requested resource could not be found"}`

	require.JSONEqual(t, got, expected)

	if got, expected := w.Code, http.StatusNotFound; got != expected {
		t.Fatalf("got status code %d, expected %d", got, expected)
	}
}

func TestMethodNotAllowedResponse(t *testing.T) {
	app := Application{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	app.MethodNotAllowedResponse(w, r)

	got := w.Body.String()
	expected := fmt.Sprintf(`{"error":"the %s method is not supported for this resource"}`, http.MethodGet)

	require.JSONEqual(t, got, expected)

	if got, expected := w.Code, http.StatusMethodNotAllowed; got != expected {
		t.Fatalf("got status code %d, expected %d", got, expected)
	}
}

func TestBadRequestResponse(t *testing.T) {
	app := Application{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	err := errors.New("error message")

	app.BadRequestResponse(w, r, err)

	got := w.Body.String()
	expected := fmt.Sprintf(`{"error":"%s"}`, err.Error())

	require.JSONEqual(t, got, expected)

	if got, expected := w.Code, http.StatusBadRequest; got != expected {
		t.Fatalf("got status code %d, expected %d", got, expected)
	}
}

func TestFailedValidationResponse(t *testing.T) {
	app := Application{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	app.FailedValidationResponse(w, r, map[string]string{"email": "wrong email"})

	got := w.Body.String()
	expected := `{"error":{"email":"wrong email"}}`

	require.JSONEqual(t, got, expected)

	if got, expected := w.Code, http.StatusUnprocessableEntity; got != expected {
		t.Fatalf("got status code %d, expected %d", got, expected)
	}
}

func TestInvalidAuthenticationTokenResponse(t *testing.T) {
	app := Application{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	app.InvalidAuthenticationTokenResponse(w, r, nil)

	got := w.Body.String()
	expected := `{"error":"invalid or missing authentication token"}`

	require.JSONEqual(t, got, expected)

	if got, expected := w.Code, http.StatusUnauthorized; got != expected {
		t.Fatalf("got status code %d, expected %d", got, expected)
	}

	if got, expected := w.Header().Get("Authorization"), "Bearer"; got != expected {
		t.Fatalf("got Authorization header with value %v, expected %v", got, expected)
	}
}

func TestAuthenticationRequiredResponse(t *testing.T) {
	app := Application{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	app.AuthenticationRequiredResponse(w, r, nil)

	got := w.Body.String()
	expected := `{"error":"you must be authenticated to access this resource"}`

	require.JSONEqual(t, got, expected)

	if got, expected := w.Code, http.StatusUnauthorized; got != expected {
		t.Fatalf("got status code %d, expected %d", got, expected)
	}
}

func TestForbiddenResponse(t *testing.T) {
	app := Application{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	app.ForbiddenResponse(w, r, nil)

	got := w.Body.String()
	expected := `{"error":"your don't have the right to access this resource"}`

	require.JSONEqual(t, got, expected)

	if got, expected := w.Code, http.StatusForbidden; got != expected {
		t.Fatalf("got status code %d, expected %d", got, expected)
	}
}

func TestRateLimitExceededResponse(t *testing.T) {
	app := Application{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	app.RateLimitExceededResponse(w, r)

	got := w.Body.String()
	expected := `{"error":"rate limit exceeded"}`

	require.JSONEqual(t, got, expected)

	if got, expected := w.Code, http.StatusTooManyRequests; got != expected {
		t.Fatalf("got status code %d, expected %d", got, expected)
	}
}
