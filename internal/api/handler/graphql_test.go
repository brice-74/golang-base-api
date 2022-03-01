package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/api/handler"
	"github.com/brice-74/golang-base-api/internal/testutils/mocks"
	"github.com/brice-74/golang-base-api/internal/testutils/require"
)

func TestGraphQL(t *testing.T) {
	req, err := http.NewRequest("GET", "/", strings.NewReader(`{"query":"{queryCheck}"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	app := &application.Application{}
	handler.GraphQL(app).ServeHTTP(rr, req)

	expected := `{"data":{"queryCheck":"ok"}}`
	require.JSONEqual(t, rr.Body.String(), expected)
}

func TestGraphQLPanic(t *testing.T) {
	req, err := http.NewRequest("GET", "/", strings.NewReader(`{"query":"{queryPanic(panic:true)}"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	app := &application.Application{
		Logger: mocks.NewLogger(),
	}

	handler.GraphQL(app).ServeHTTP(rr, req)

	expected := `{"errors":[{"message":"panic occurred: I panic !!!","path":["queryPanic"]}],"data":null}`
	require.JSONEqual(t, rr.Body.String(), expected)
}
