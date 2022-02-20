package main

import (
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/pkg/jsonlog"
)

func main() {
	cfg := getConfigFromFlags()

	logger := jsonlog.New(
		os.Stdout,
		jsonlog.LevelInfo,
		jsonlog.Middlewares{
			AfterPrintError: func(err error) {
				sentry.CaptureException(err)
			},
		},
	)

	postgres, err := openPostgresDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer postgres.Close()

	logger.PrintInfo("postgres connection pool established", nil)

	m := application.NewModels(postgres)

	app := &application.Application{
		Config: cfg,
		Models: m,
		Logger: logger,
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn:         app.Config.Sentry.DSN,
		Environment: app.Config.Env,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	err = serve(app)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
