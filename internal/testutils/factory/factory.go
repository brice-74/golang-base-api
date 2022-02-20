package factory

import (
	"database/sql"
	"testing"

	"github.com/jaswdr/faker"
)

type Factory struct {
	T  *testing.T
	DB *sql.DB

	faker *faker.Faker
}

func New(t *testing.T, db *sql.DB) *Factory {
	fak := faker.New()

	return &Factory{
		T:     t,
		DB:    db,
		faker: &fak,
	}
}
