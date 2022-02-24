package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/api/handler"
	"github.com/brice-74/golang-base-api/internal/testutils/require"
)

func TestHealthcheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	app := &application.Application{}
	app.Config.Env = "dev"

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler.Healthcheck(app))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %d want %d", status, http.StatusOK)
	}

	expected := `{"status":"available","systemInfo":{"environment":"dev"}}`
	require.JSONEqual(t, rr.Body.String(), expected)
}
