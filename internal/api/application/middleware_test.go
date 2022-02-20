package application_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/google/go-cmp/cmp"
)

func TestEnableCORS(t *testing.T) {
	tests := []struct {
		title           string
		origin          string
		headers         map[string]string
		expectedHeaders map[string]string
	}{
		{
			title: "should not allow origin",
		},
		{
			title:  "should allow trusted header",
			origin: "http://testing",
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "http://testing",
			},
		},
		{
			title:  "should return headers for option request",
			origin: "http://testing",
			headers: map[string]string{
				"Access-Control-Request-Method": "POST",
			},
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "http://testing",
				"Access-Control-Allow-Methods": "OPTIONS, PUT, PATCH, DELETE",
				"Access-Control-Allow-Headers": "Authorization, Content-Type",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			app := application.Application{
				Config: application.Config{
					CORS: struct{ TrustedOrigins []string }{
						TrustedOrigins: []string{tt.origin},
					},
				},
			}

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if diff := cmp.Diff(
					w.Header().Values("Vary"),
					[]string{"Origin", "Access-Control-Request-Method"},
				); diff != "" {
					t.Error(diff)
				}

				if len(tt.expectedHeaders) > 0 {
					for k, v := range tt.expectedHeaders {
						if val := w.Header().Get(k); val != v {
							t.Errorf("got header value %s with value %s, expected %s", k, val, v)
						}
					}
				}
			})

			handler := app.EnableCORS(next)

			req := httptest.NewRequest("OPTIONS", "http://testing", nil)
			req.Header.Add("Origin", tt.origin)

			for k, v := range tt.headers {
				req.Header.Add(k, v)
			}

			handler.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
}
