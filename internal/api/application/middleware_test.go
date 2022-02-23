package application_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/pkg/jsonlog"
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

func TestRecoverPanic(t *testing.T) {
	app := &application.Application{
		Logger: jsonlog.New(
			os.Stdout,
			// Disable logs
			jsonlog.LevelDisable,
			jsonlog.Middlewares{},
		),
	}

	handlerFunc := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("I panic !!!")
	})

	req := httptest.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	app.RecoverPanic(handlerFunc).ServeHTTP(res, req)

	if res.Header().Get("Connection") != "close" {
		t.Errorf("Connection header should be close")
	}

	if res.Code != 500 {
		t.Errorf("Response http code should be 500")
	}

	expect := `{"error":"the server encountered a problem and could not process your request"}`
	if strings.TrimSpace(res.Body.String()) != expect {
		t.Fatalf("handler returned unexpected body: got %s want %s", res.Body.String(), expect)
	}
}

func TestRateLimit(t *testing.T) {
	app := &application.Application{}
	app.Config.Limiter.Enabled = true
	app.Config.Limiter.RPS = 2
	app.Config.Limiter.Burst = 4

	req := httptest.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	limiter := app.RateLimit(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))

	for i := 0; i < app.Config.Limiter.Burst+1; i++ {
		limiter.ServeHTTP(res, req)

		switch {
		case i == app.Config.Limiter.Burst:
			if res.Code != 429 {
				t.Fatalf("response code http must be 429 at request %d, got %d", i, res.Code)
			}

			expect := `{"error":"rate limit exceeded"}`
			if strings.TrimSpace(res.Body.String()) != expect {
				t.Fatalf("limiter returned unexpected body: got %s want %s", res.Body.String(), expect)
			}
		default:
			if res.Code != 200 {
				t.Errorf("response code http must be 200 at request %d, got %d", i, res.Code)
			}
		}
	}
}
