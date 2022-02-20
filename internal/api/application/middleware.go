package application

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/brice-74/golang-base-api/internal/domains/user"
)

// Allow CORS for specific domains.
func (app *Application) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		if origin != "" && len(app.Config.CORS.TrustedOrigins) != 0 {
			for i := range app.Config.CORS.TrustedOrigins {
				if origin == app.Config.CORS.TrustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					// Handle Preflight requests.
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
						return
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// RecoverPanic sends a 500 server error instead of just closing the HTTP connection.
func (app *Application) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Use the builtin recover function to check if there has been a panic or // not.
			if err := recover(); err != nil {
				// If there was a panic, set a "Connection: close" header on the response.
				// This acts as a trigger to make Go's HTTP server automatically close the
				// current connection after a response has been // sent.
				w.Header().Set("Connection", "close")

				app.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// RateLimit adds rate limiters for each client IP addresses.
func (app *Application) RateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Background goroutine which removes old entries from the clients map once every minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			// Loop through all clients. If they haven't been seen within the last three
			// minutes, delete the corresponding entry from the map.
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.Config.Limiter.Enabled {
			// Retrieve IP from the client.
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.ServerErrorResponse(w, r, err)
				return
			}

			mu.Lock()

			if _, found := clients[ip]; !found {
				// 2 seconds with a maximum of 4 requests in a single burst.
				clients[ip] = &client{
					limiter: rate.NewLimiter(
						rate.Limit(app.Config.Limiter.RPS),
						app.Config.Limiter.Burst,
					),
				}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.RateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}

// Authenticate gets the token from the Authorization header and adds the retrieved user to the HTTP context request.
func (app *Application) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Indicates to any cache systems that the response can vary based on the Authorization header.
		w.Header().Add("Vary", "Authorization")

		var authorizationHeader = r.Header.Get("Authorization")
		var grantTypeHeader = r.Header.Get("Grant-Type")
		// Set an anonymous user in the request is no Authorization header.
		if authorizationHeader == "" {
			ctx := app.ContextWithUser(r.Context(), &UserCtx{User: user.AnonymousUser})
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		// Split Authorization header to recover token
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.InvalidAuthenticationTokenResponse(w, r)
			return
		}
		// Check Grant-Type header to know if is a resfresh token
		var isAccess bool = true
		var secret string = app.Config.JWT.Access.Secret
		if grantTypeHeader != "" {
			if grantTypeHeader == "refresh-token" {
				isAccess = false
				secret = app.Config.JWT.Refresh.Secret
			} else {
				w.Header().Set("Grant-Type", "refresh-token")
				app.BadRequestResponse(w, r, errors.New("Invalid header value"))
				return
			}
		}
		// Check token is valid and up to date
		token, err := VerifyToken(headerParts[1], secret)
		if err != nil {
			app.InvalidAuthenticationTokenResponseMsg(w, r, err.Error())
			return
		}
		// Extract userID claim in the token
		claims, err := ExtractTokenMetadata(token, []JwtClaimKey{UserIdClaim, UserAgentIdClaim})
		if err != nil {
			app.InvalidAuthenticationTokenResponseMsg(w, r, err.Error())
			return
		}

		u, err := app.Models.User.GetById(claims[UserIdClaim])
		if err != nil {
			switch {
			case errors.Is(err, user.ErrNotFound):
				app.NotFoundResponseMsg(w, r, err.Error())
			default:
				app.ServerErrorResponse(w, r, err)
			}
		}

		ctx := app.ContextWithUser(r.Context(), &UserCtx{
			User: u,
			Token: &TokenCtx{
				IsAccess: isAccess,
			},
			SessionID: claims[UserAgentIdClaim],
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
