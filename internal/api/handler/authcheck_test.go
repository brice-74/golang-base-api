package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/api/handler"
	"github.com/brice-74/golang-base-api/internal/domains/user"
)

func TestAuthToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	app := &application.Application{}

	ctx := app.ContextWithUser(req.Context(), &application.UserCtx{
		User: &user.User{Roles: user.Roles{user.RoleUser}},
		Client: &application.Client{
			SessionID: "1234",
			IP:        "0.0.0.0",
			Agent:     "agent",
		},
	})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler.AuthToken(app))

	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"Client":{"Agent":"agent","IP":"0.0.0.0","Session":"1234"},"Roles":["ROLE_USER"]}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Fatalf("handler returned unexpected body: got %s want %s", rr.Body.String(), expected)
	}
}
