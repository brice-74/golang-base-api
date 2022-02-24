package handler

import (
	"net/http"

	"github.com/brice-74/golang-base-api/internal/api/application"
)

func AuthToken(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.ClientFromContext(r.Context())

		context := application.Envelope{
			"Client": application.Envelope{
				"Session": ctx.Session.ID,
				"Agent":   ctx.Agent.Agent,
				"IP":      ctx.Agent.IP,
			},
			"Roles": ctx.User.Roles,
		}

		if err := app.WriteJSON(w, http.StatusOK, context, nil); err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}
