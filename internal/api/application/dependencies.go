package application

import (
	"database/sql"

	"github.com/brice-74/golang-base-api/internal/domains/user"
)

type Models struct {
	User user.Model
}

func NewModels(db *sql.DB) Models {
	return Models{
		User: user.Model{DB: db},
	}
}
