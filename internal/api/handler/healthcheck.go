package handler

import (
	"net/http"

	"github.com/brice-74/golang-base-api/internal/api/application"
)

// Healthcheck is an endpoint checking if the health of the API
func Healthcheck(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		env := application.Envelope{
			"status": "available",
			"systemInfo": map[string]string{
				"environment": app.Config.Env,
			},
		}

		if err := app.WriteJSON(w, http.StatusOK, env, nil); err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}
