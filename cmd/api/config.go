package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/brice-74/golang-base-api/internal/api/application"
)

func getConfigFromFlags() application.Config {
	var cfg application.Config

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(fmt.Errorf("error when parsing port: %w", err))
	}

	if port == 0 {
		port = 4000
	}

	flag.IntVar(&cfg.Port, "port", port, "API server port")
	flag.StringVar(&cfg.Env, "env", os.Getenv("ENV"), "Environment (dev|staging|prod)")

	// Database
	flag.StringVar(&cfg.DB.URL, "db-url", os.Getenv("DATABASE_URL"), "PostgreSQL URL")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	// Rate limiter configuration
	flag.Float64Var(&cfg.Limiter.RPS, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	// CORS trusted domains
	var trustedOrigins string
	flag.StringVar(
		&trustedOrigins,
		"cors-trusted-origins",
		os.Getenv("CORS_TRUSTED_ORIGINS"),
		"Trusts CORS domains origins (space separated)",
	)

	// Integrations
	flag.StringVar(&cfg.Sentry.DSN, "sentry-dsn", os.Getenv("SENTRY_DSN"), "DSN for Sentry integrations")

	// JWT
	flag.StringVar(&cfg.JWT.Access.Secret, "jwt-access-secret", os.Getenv("JWT_ACCESS_SECRET"), "Secret key using to secure access JWT")
	flag.StringVar(&cfg.JWT.Access.Expiration, "jwt-access-expiration-time", "15m", "Validity time of access JWT")
	flag.StringVar(&cfg.JWT.Refresh.Secret, "jwt-refresh-secret", os.Getenv("JWT_REFRESH_SECRET"), "Secret key using to secure refresh JWT")
	flag.StringVar(&cfg.JWT.Refresh.Expiration, "jwt-refresh-expiration-time", "168h", "Validity time of refresh JWT")

	flag.Parse()

	cfg.CORS.TrustedOrigins = strings.Fields(trustedOrigins)

	return cfg
}
