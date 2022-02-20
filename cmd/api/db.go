package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/brice-74/golang-base-api/internal/api/application"
)

func openPostgresDB(cfg application.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.URL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Establish a new connection to the database.
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
