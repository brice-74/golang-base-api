package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/api/handler"
)

func Routes(app *application.Application) http.Handler {
	router := httprouter.New()
	// Customize default router error responses.
	router.NotFound = http.HandlerFunc(app.NotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.MethodNotAllowedResponse)

	//----------------//
	//			REST			//
	//----------------//

	if app.Config.Env == "dev" {
		router.HandlerFunc(http.MethodGet, "/check/health", handler.Healthcheck(app))
		router.HandlerFunc(http.MethodGet, "/check/token", handler.AuthToken(app))
	}

	//-------------------//
	//			GraphQL			 //
	//-------------------//

	router.HandlerFunc(http.MethodPost, "/graphql", handler.GraphQL(app))

	return app.RecoverPanic(app.EnableCORS(app.RateLimit(app.Authenticate(router))))
}
