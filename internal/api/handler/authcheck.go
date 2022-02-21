package handler

import (
	"net/http"

	"github.com/brice-74/golang-base-api/internal/api/application"
)

func AuthToken(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.UserFromContext(r.Context())

		context := application.Envelope{
			"session": ctx.SessionID,
			"roles":   ctx.User.Roles,
		}

		if err := app.WriteJSON(w, http.StatusOK, context, nil); err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}
