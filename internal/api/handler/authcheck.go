package handler

import (
	"net/http"

	"github.com/brice-74/golang-base-api/internal/api/application"
)

func AuthToken(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actx := app.UserFromContext(r.Context())

		var scope string
		if actx.Token.IsAccess {
			scope = "access token"
		} else {
			scope = "refresh token"
		}

		context := application.Envelope{
			"session": actx.SessionID,
			"roles":   actx.User.Roles,
			"token": application.Envelope{
				"scope": scope,
			},
		}

		if err := app.WriteJSON(w, http.StatusOK, context, nil); err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}
