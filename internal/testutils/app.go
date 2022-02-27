package testutils

import (
	"database/sql"

	"github.com/brice-74/golang-base-api/internal/api/application"
)

func NewApplication(db *sql.DB) *application.Application {
	app := &application.Application{
		Models: application.NewModels(db),
	}
	app.Config.JWT.Access.Secret = "secret access"
	app.Config.JWT.Access.Expiration = "3m"
	app.Config.JWT.Refresh.Secret = "secret refresh"
	app.Config.JWT.Refresh.Expiration = "10m"

	return app
}
