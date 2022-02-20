package application

import (
	"errors"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestCreateTokens(t *testing.T) {
	app := &Application{}
	app.Config.JWT.Access.Secret = "access-secret"
	app.Config.JWT.Access.Expiration = "15m"
	app.Config.JWT.Refresh.Secret = "refresh-secret"
	app.Config.JWT.Refresh.Expiration = "15m"

	const (
		userID    string = "1234"
		sessionID string = "5678"
	)

	td, err := app.CreateTokens(userID, sessionID)
	if err != nil {
		t.Fatal(err)
	}

	at, err := jwt.Parse(td.AccessToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(app.Config.JWT.Access.Secret), nil
	})
	if err != nil {
		t.Fatal(err)
	}

	rt, err := jwt.Parse(td.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(app.Config.JWT.Refresh.Secret), nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if !at.Valid {
		t.Fatal("access token is invalid")
	}
	if !rt.Valid {
		t.Fatal("refresh token is invalid")
	}
}

func TestExtractTokenMetadata(t *testing.T) {
	const (
		testKey1 JwtClaimKey = "keytest1"
		valKey1  string      = "val1"
	)
	var claims = jwt.MapClaims{}
	claims[string(testKey1)] = valKey1

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	claimMap, err := ExtractTokenMetadata(token, []JwtClaimKey{
		testKey1,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, expect := claimMap[testKey1], valKey1; got != expect {
		t.Fatalf("got claim value %s, expected %s", got, expect)
	}
}

func TestVerifyToken(t *testing.T) {
	const secret string = "secret"

	var claims = jwt.MapClaims{}
	claims["exp"] = time.Time{}.Unix()

	expiresToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}

	badToken, err := jwt.New(jwt.SigningMethodHS384).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}

	goodToken, err := jwt.New(jwt.SigningMethodHS256).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Should be expire error", func(t *testing.T) {
		_, err := VerifyToken(expiresToken, secret)
		expect := errors.New("Token is expired")
		if err == nil {
			t.Fatal("nil error")
		}
		if err.Error() != expect.Error() {
			t.Fatalf("Unexpected error got: %s, but expected: %s", err.Error(), expect.Error())
		}
	})

	t.Run("Should be signing method error", func(t *testing.T) {
		_, err := VerifyToken(badToken, secret)
		expect := errors.New("unexpected signing method: HS384")
		if err == nil {
			t.Fatal("nil error")
		}
		if err.Error() != expect.Error() {
			t.Fatalf("Unexpected error got: %s, but expected: %s", err.Error(), expect.Error())
		}
	})

	t.Run("Should be good token", func(t *testing.T) {
		_, err := VerifyToken(goodToken, secret)
		if err != nil {
			t.Fatalf("Token should be good but error occured: %s", err.Error())
		}
	})
}
