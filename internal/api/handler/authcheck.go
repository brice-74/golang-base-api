package handler

import (
	"net/http"

	"github.com/brice-74/golang-base-api/internal/api/application"
)

func AuthToken(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.UserFromContext(r.Context())

		context := application.Envelope{
			"Client": application.Envelope{
				"Session": ctx.Client.SessionID,
				"Agent":   ctx.Client.Agent,
				"IP":      ctx.Client.IP,
			},
			"Roles": ctx.User.Roles,
		}

		if err := app.WriteJSON(w, http.StatusOK, context, nil); err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}
