package mocks

import (
	"testing"
	"time"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/dgrijalva/jwt-go"
)

func CreateToken(t *testing.T, method *jwt.SigningMethodHMAC, claims jwt.MapClaims, secret string) string {
	to := jwt.NewWithClaims(method, claims)
	toStr, err := to.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("error during jwt token creation: %s", err.Error())
	}
	return toStr
}

func CreateClaims(userID, sessionID string, time time.Time) jwt.MapClaims {
	return jwt.MapClaims{
		string(application.SessionIdClaim): sessionID,
		string(application.UserIdClaim):    userID,
		string(application.ExpireClaim):    time.Unix(),
	}
}
