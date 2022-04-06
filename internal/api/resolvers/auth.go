package resolvers

import (
	"context"
	"errors"
	"time"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/domains/user"
	"github.com/brice-74/golang-base-api/pkg/validator"
	"github.com/graph-gophers/graphql-go"
	"github.com/ventu-io/go-shortid"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUserAccount: register a new user account
func (r Root) RegisterUserAccount(_ context.Context, params RegisterUserAccountParams) (*UserAccountResolver, error) {
	u := user.User{
		Email:      params.Input.Email,
		Password:   params.Input.Password,
		ProfilName: params.Input.ProfilName,
	}
	// check that all entries are valid
	v := validator.New()
	u.ValidateEmailEntry(v)
	u.ValidatePasswordEntry(v)
	u.ValidateProfilNameEntry(v)
	if !v.Valid() {
		return nil, validatorError{Errors: v.Errors}
	}
	// generate short identifier
	if short, err := shortid.Generate(); err != nil {
		return nil, err
	} else {
		u.ShortId = short
	}
	// hash password
	if hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14); err != nil {
		return nil, err
	} else {
		u.Password = string(hash)
	}
	// add user role
	u.Roles = []user.Role{user.RoleUser}
	// insert peacefully
	if err := r.App.Models.User.InsertRegisteredUserAccount(&u); err != nil {
		return nil, resolverErrDatabaseOperation(err)
	}

	return &UserAccountResolver{app: r.App, user: u}, nil
}

type RegisterUserAccountParams struct {
	Input RegisterUserAccountInput
}

type RegisterUserAccountInput struct {
	Email      string
	Password   string
	ProfilName string
}

// LoginUserAccount: authenticate a user by returning tokens
func (r Root) LoginUserAccount(ctx context.Context, params LoginUserAccountParams) (*TokensUserAccountResolver, error) {
	uctx := r.App.ClientFromContext(ctx)

	uEntry := user.User{
		Email:    params.Email,
		Password: params.Password,
	}
	// check that all entries are valid
	v := validator.New()
	uEntry.ValidateEmailEntry(v)
	uEntry.ValidatePasswordEntry(v)
	if !v.Valid() {
		return nil, validatorError{Errors: v.Errors}
	}
	// find registered user
	uReg, err := r.App.Models.User.GetByEmail(uEntry.Email)
	if err != nil {
		if err == user.ErrNotFoundUser {
			return nil, resolverErrNotFound(err)
		} else {
			return nil, resolverErrDatabaseOperation(err)
		}
	}
	// check password
	if err = bcrypt.CompareHashAndPassword([]byte(uReg.Password), []byte(uEntry.Password)); err != nil {
		return nil, resolverErrUnauthorized(errors.New("incorrect password"))
	}
	// create jwt access & refresh
	td, err := r.App.CreateTokens(uReg.ID, string(params.SessionID))
	if err != nil {
		return nil, err
	}
	// insert session information or update if user need re login
	if err = r.App.Models.User.InsertOrUpdateUserSession(
		&user.Session{
			ID:            string(params.SessionID),
			DeactivatedAt: time.Unix(td.RefreshExp, 0),
			IP:            uctx.Agent.IP,
			Agent:         uctx.Agent.Agent,
			UserID:        uReg.ID,
		},
	); err != nil {
		return nil, resolverErrDatabaseOperation(err)
	}
	// everything is good, return tokens using resolver
	return &TokensUserAccountResolver{app: r.App, tokens: user.Tokens{
		Access:  td.AccessToken,
		Refresh: td.RefreshToken,
	}}, nil
}

type LoginUserAccountParams struct {
	Email     string
	Password  string
	SessionID graphql.ID
}

func (r Root) RefreshUserAccount(ctx context.Context, params RefreshUserAccountParams) (*TokensUserAccountResolver, error) {
	uctx := r.App.ClientFromContext(ctx)
	// check token is valid and up to date
	token, err := application.VerifyToken(params.Token, r.App.Config.JWT.Refresh.Secret)
	if err != nil {
		return nil, err
	}
	// extract token claims
	claims, err := application.ExtractTokenMetadata(token, []application.JwtClaimKey{application.UserIdClaim, application.SessionIdClaim})
	if err != nil {
		return nil, errors.New("Required claims from token not found")
	}
	// get session and verify that user id claim is associated to session id claim
	_, s, err := r.App.Models.User.GetUserAndSession(claims[application.UserIdClaim], claims[application.SessionIdClaim])
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFoundUserAndSession):
			return nil, resolverErrNotFound(err)
		default:
			return nil, resolverErrDatabaseOperation(err)
		}
	}
	// create new tokens
	td, err := r.App.CreateTokens(claims[application.UserIdClaim], claims[application.SessionIdClaim])
	if err != nil {
		return nil, err
	}
	// update session information
	if err = r.App.Models.User.InsertOrUpdateUserSession(
		&user.Session{
			ID:            s.ID,
			DeactivatedAt: time.Unix(td.RefreshExp, 0),
			IP:            uctx.Agent.IP,
			Agent:         uctx.Agent.Agent,
			UserID:        s.UserID,
		},
	); err != nil {
		return nil, resolverErrDatabaseOperation(err)
	}

	return &TokensUserAccountResolver{app: r.App, tokens: user.Tokens{
		Access:  td.AccessToken,
		Refresh: td.RefreshToken,
	}}, nil
}

type RefreshUserAccountParams struct {
	Token string
}

type TokensUserAccountResolver struct {
	app    *application.Application
	tokens user.Tokens
}

func (r TokensUserAccountResolver) Access() string {
	return r.tokens.Access
}

func (r TokensUserAccountResolver) Refresh() string {
	return r.tokens.Refresh
}

func (r Root) LogoutUserAccount(ctx context.Context) (bool, error) {
	c := r.App.ClientFromContext(ctx)

	if err := r.App.Models.User.InsertOrUpdateUserSession(
		&user.Session{
			ID:            c.Session.ID,
			DeactivatedAt: time.Now(),
			IP:            c.Agent.IP,
			Agent:         c.Agent.Agent,
			UserID:        c.User.ID,
		},
	); err != nil {
		return false, resolverErrDatabaseOperation(err)
	}

	return true, nil
}
