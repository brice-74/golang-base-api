package application

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthError error

var (
	InvalidToken AuthError = errors.New("Invalid token")
)

type JwtClaimKey string

const (
	UserIdClaim    JwtClaimKey = "user_id"
	SessionIdClaim JwtClaimKey = "user_agent_id"
	ExpireClaim    JwtClaimKey = "exp"
)

type TokensDetails struct {
	AccessToken  string
	RefreshToken string
	RefreshExp   int64
}

func (app *Application) CreateTokens(userID string, userAgentID string) (*TokensDetails, error) {
	var td = &TokensDetails{}

	accessDuration, err := time.ParseDuration(app.Config.JWT.Access.Expiration)
	if err != nil {
		return nil, err
	}

	refreshDuration, err := time.ParseDuration(app.Config.JWT.Refresh.Expiration)
	if err != nil {
		return nil, err
	}
	td.RefreshExp = time.Now().Add(refreshDuration).Unix()

	var atClaims = jwt.MapClaims{}
	atClaims[string(SessionIdClaim)] = userAgentID
	atClaims[string(UserIdClaim)] = userID
	atClaims[string(ExpireClaim)] = time.Now().Add(accessDuration).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(app.Config.JWT.Access.Secret))
	if err != nil {
		return nil, err
	}

	var rtClaims = jwt.MapClaims{}
	rtClaims[string(SessionIdClaim)] = userAgentID
	rtClaims[string(UserIdClaim)] = userID
	rtClaims[string(ExpireClaim)] = td.RefreshExp
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(app.Config.JWT.Refresh.Secret))
	if err != nil {
		return nil, err
	}
	return td, nil
}

// Make sure that the token method conform to "SigningMethodHMAC" and is up to date
func VerifyToken(bearer string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(string(bearer), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, InvalidToken
	}
	return token, nil
}

func ExtractTokenMetadata(token *jwt.Token, claimKeys []JwtClaimKey) (map[JwtClaimKey]string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		var returned = make(map[JwtClaimKey]string)
		for _, claimKey := range claimKeys {
			claim, ok := claims[string(claimKey)].(string)
			if !ok {
				return nil, fmt.Errorf("Jwt claim key not found: %s", claimKey)
			}
			returned[claimKey] = claim
		}
		return returned, nil
	}
	return nil, InvalidToken
}
