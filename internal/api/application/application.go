package application

import (
	"github.com/brice-74/golang-base-api/pkg/jsonlog"
)

type Application struct {
	Config Config
	Models Models
	Logger jsonlog.Logger
}

type Config struct {
	Port int
	Env  string
	DB   struct {
		URL          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Limiter struct {
		RPS     float64
		Burst   int
		Enabled bool
	}
	CORS struct {
		TrustedOrigins []string
	}
	Sentry struct {
		DSN string
	}
	JWT struct {
		Access struct {
			Secret     string
			Expiration string
		}
		Refresh struct {
			Secret     string
			Expiration string
		}
	}
}

type Envelope map[string]interface{}
