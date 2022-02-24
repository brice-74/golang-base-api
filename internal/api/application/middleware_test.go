package application_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/domains/user"
	"github.com/brice-74/golang-base-api/internal/testutils"
	"github.com/brice-74/golang-base-api/internal/testutils/factory"
	"github.com/brice-74/golang-base-api/pkg/jsonlog"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-cmp/cmp"
	"github.com/twinj/uuid"
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

func TestAuthenticate(t *testing.T) {
	var (
		db  = testutils.PrepareDB(t)
		fac = factory.New(t, db)
		app = &application.Application{
			Logger: jsonlog.New(
				os.Stdout,
				jsonlog.LevelDisable,
				jsonlog.Middlewares{},
			),
			Models: application.NewModels(db),
		}
	)

	app.Config.JWT.Access.Secret = "secret"

	var createToken = func(t *testing.T, method *jwt.SigningMethodHMAC, claims jwt.MapClaims, secret string) string {
		to := jwt.NewWithClaims(method, claims)
		toStr, err := to.SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("error during jwt token creation: %s", err.Error())
		}
		return toStr
	}

	u := fac.CreateUserAccount(nil)
	s := fac.CreateUserSession(&user.Session{
		UserID: u.ID,
	})

	goodClaims := jwt.MapClaims{
		string(application.SessionIdClaim): s.ID,
		string(application.UserIdClaim):    u.ID,
		string(application.ExpireClaim):    time.Now().Add(time.Minute * 3).Unix(),
	}

	expClaims := jwt.MapClaims{
		string(application.SessionIdClaim): s.ID,
		string(application.UserIdClaim):    u.ID,
		string(application.ExpireClaim):    time.Time{}.Unix(),
	}

	randomIdsClaims := jwt.MapClaims{
		string(application.SessionIdClaim): uuid.NewV4().String(),
		string(application.UserIdClaim):    uuid.NewV4().String(),
		string(application.ExpireClaim):    time.Now().Add(time.Minute * 3).Unix(),
	}

	badIdsClaims := jwt.MapClaims{
		string(application.SessionIdClaim): "",
		string(application.UserIdClaim):    "",
		string(application.ExpireClaim):    time.Now().Add(time.Minute * 3).Unix(),
	}

	a := &application.Agent{
		IP:    "0.0.0.0",
		Agent: "agent",
	}

	type ExpectError struct {
		code int
		json string
	}

	tests := []struct {
		title           string
		headers         map[string]string
		expectedHeaders map[string]string
		expectContext   *application.ClientCtx
		expectError     *ExpectError
	}{
		{
			title: "should return context with session",
			headers: map[string]string{
				"Authorization": "Bearer " + createToken(t, jwt.SigningMethodHS256, goodClaims, app.Config.JWT.Access.Secret),
			},
			expectedHeaders: map[string]string{
				"Vary": "Authorization",
			},
			expectContext: &application.ClientCtx{
				User:    u,
				Session: s,
				Agent:   a,
			},
		},
		{
			title: "should return context with anonymous user",
			expectedHeaders: map[string]string{
				"Vary": "Authorization",
			},
			expectContext: &application.ClientCtx{
				User:  user.AnonymousUser,
				Agent: a,
			},
		},
		{
			title: "should return invalid or missing authentication token",
			headers: map[string]string{
				"Authorization": "bad bearer",
			},
			expectedHeaders: map[string]string{
				"Vary":          "Authorization",
				"Authorization": "Bearer",
			},
			expectError: &ExpectError{
				code: 401,
				json: `{"error":"invalid or missing authentication token"}`,
			},
		},
		{
			title: "should return unexpected signing method",
			headers: map[string]string{
				"Authorization": "Bearer " + createToken(t, jwt.SigningMethodHS384, goodClaims, app.Config.JWT.Access.Secret),
			},
			expectedHeaders: map[string]string{
				"Vary":          "Authorization",
				"Authorization": "Bearer",
			},
			expectError: &ExpectError{
				code: 401,
				json: `{"error":"unexpected signing method: HS384"}`,
			},
		},
		{
			title: "should return expired token",
			headers: map[string]string{
				"Authorization": "Bearer " + createToken(t, jwt.SigningMethodHS256, expClaims, app.Config.JWT.Access.Secret),
			},
			expectedHeaders: map[string]string{
				"Vary":          "Authorization",
				"Authorization": "Bearer",
			},
			expectError: &ExpectError{
				code: 401,
				json: `{"error":"Token is expired"}`,
			},
		},
		{
			title: "should return invalid signature token",
			headers: map[string]string{
				"Authorization": "Bearer " + createToken(t, jwt.SigningMethodHS256, nil, "bad secret"),
			},
			expectedHeaders: map[string]string{
				"Vary":          "Authorization",
				"Authorization": "Bearer",
			},
			expectError: &ExpectError{
				code: 401,
				json: `{"error":"signature is invalid"}`,
			},
		},
		{
			title: "should return claims not found from token",
			headers: map[string]string{
				"Authorization": "Bearer " + createToken(t, jwt.SigningMethodHS256, nil, app.Config.JWT.Access.Secret),
			},
			expectedHeaders: map[string]string{
				"Vary":          "Authorization",
				"Authorization": "Bearer",
			},
			expectError: &ExpectError{
				code: 401,
				json: `{"error":"Required claims from token not found"}`,
			},
		},
		{
			title: "should return not found",
			headers: map[string]string{
				"Authorization": "Bearer " + createToken(t, jwt.SigningMethodHS256, randomIdsClaims, app.Config.JWT.Access.Secret),
			},
			expectedHeaders: map[string]string{
				"Vary": "Authorization",
			},
			expectError: &ExpectError{
				code: 404,
				json: `{"error":"User or user session not found"}`,
			},
		},
		{
			title: "should return server error",
			headers: map[string]string{
				"Authorization": "Bearer " + createToken(t, jwt.SigningMethodHS256, badIdsClaims, app.Config.JWT.Access.Secret),
			},
			expectedHeaders: map[string]string{
				"Vary": "Authorization",
			},
			expectError: &ExpectError{
				code: 500,
				json: `{"error":"the server encountered a problem and could not process your request"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotCtx := app.ClientFromContext(r.Context())

				if diff := cmp.Diff(tt.expectContext, gotCtx); diff != "" {
					t.Errorf(diff)
				}
			})

			req := httptest.NewRequest("GET", "/", nil)
			// fix request addr and agent
			req.RemoteAddr = a.IP
			req.Header.Set("User-Agent", a.Agent)

			for k, v := range tt.headers {
				req.Header.Add(k, v)
			}

			res := httptest.NewRecorder()
			app.Authenticate(next).ServeHTTP(res, req)

			if len(tt.expectedHeaders) > 0 {
				for k, v := range tt.expectedHeaders {
					if got := res.Header().Get(k); got != v {
						t.Errorf("got header value %s with value %s, expected %s", k, got, v)
					}
				}
			}

			if tt.expectError != nil {
				resErr := strings.TrimSpace(res.Body.String())
				if tt.expectError.code != res.Code || tt.expectError.json != resErr {
					t.Errorf("got code: %d and error: %s \n expect code: %d and error: %s",
						res.Code,
						resErr,
						tt.expectError.code,
						tt.expectError.json,
					)
				}
			}
		})
	}
}
