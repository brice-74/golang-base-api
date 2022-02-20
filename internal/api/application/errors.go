package application

import (
	"fmt"
	"net/http"
)

// ErrorResponse is a generic HTTP error helper.
func (app *Application) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := Envelope{"error": message}

	err := app.WriteJSON(w, status, env, nil)
	if err != nil {
		app.LogError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// ServerErrorResponse returns a 500 error to the client.
func (app *Application) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.LogError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

// NotFoundResponse returns a 404 error to the client.
func (app *Application) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.ErrorResponse(w, r, http.StatusNotFound, message)
}

// NotFoundResponse returns a 404 error to the client with custom message.
func (app *Application) NotFoundResponseMsg(w http.ResponseWriter, r *http.Request, message string) {
	app.ErrorResponse(w, r, http.StatusNotFound, message)
}

// MethodNotAllowedResponse returns a 405 error to the client.
func (app *Application) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// BadRequestResponse returns a 400 error to the client.
func (app *Application) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

// FailedValidationResponse returns a 422 error to the client due to validator errors.
func (app *Application) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors interface{}) {
	app.ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// InvalidAuthenticationTokenResponse returns a 401 error indicating the token is not valid.
func (app *Application) InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	// Indicates to the client we expect a bearer token.
	w.Header().Set("Authorization", "Bearer")

	message := "invalid or missing authentication token"
	app.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

// InvalidAuthenticationTokenResponse returns a 401 error indicating the token is not valid with custom message.
func (app *Application) InvalidAuthenticationTokenResponseMsg(w http.ResponseWriter, r *http.Request, message string) {
	// Indicates to the client we expect a bearer token.
	w.Header().Set("Authorization", "Bearer")
	app.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

// AuthenticationRequiredResponse returns a 401 response to the client.
func (app *Application) AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

// ForbiddenResponse returns a 403 response to the client.
func (app *Application) ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	message := "your don't have the right to access this resource"
	app.ErrorResponse(w, r, http.StatusForbidden, message)
}

// ForbiddenResponse returns a 403 response to the client with custom message.
func (app *Application) ForbiddenResponseMsg(w http.ResponseWriter, r *http.Request, message string) {
	app.ErrorResponse(w, r, http.StatusForbidden, message)
}

// RateLimitExceededResponse returns a 429 response to the client.
func (app *Application) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.ErrorResponse(w, r, http.StatusTooManyRequests, message)
}
