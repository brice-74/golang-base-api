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

type UserAccountResolver struct {
	app  *application.Application
	user user.User
}

func (r UserAccountResolver) ID() graphql.ID {
	return graphql.ID(r.user.ID)
}

func (r UserAccountResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.user.CreatedAt}
}

func (r UserAccountResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: r.user.UpdatedAt}
}

func (r UserAccountResolver) Active() bool {
	return r.user.DeactivatedAt.IsZero()
}

func (r UserAccountResolver) Email() string {
	return r.user.Email
}

func (r UserAccountResolver) Password() string {
	return r.user.Password
}

func (r UserAccountResolver) Roles() user.Roles {
	return r.user.Roles
}

func (r UserAccountResolver) ProfilName() string {
	return r.user.ProfilName
}

func (r UserAccountResolver) ShortId() string {
	return r.user.ShortId
}

// LoginUserAccount: authenticate a user by returning tokens
func (r Root) LoginUserAccount(_ context.Context, params LoginUserAccountParams) (*TokensUserAccountResolver, error) {
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
		if err == user.ErrNotFound {
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
	td, err := r.App.CreateTokens(uReg.ID, string(params.Agent.ID))
	if err != nil {
		return nil, err
	}

	if err = r.App.Models.User.InsertOrUpdateUserSession(
		&user.Session{
			ID:            string(params.Agent.ID),
			DeactivatedAt: time.Unix(td.RefreshExp, 0),
			IP:            params.Agent.IP,
			Name:          params.Agent.Name,
			Location:      params.Agent.Location,
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
	Email    string
	Password string
	Agent    AgentParams
}

type AgentParams struct {
	ID       graphql.ID
	IP       string
	Name     string
	Location string
}

func (r Root) RefreshUserAccount(_ context.Context, params RefreshUserAccountParams) (*TokensUserAccountResolver, error) {
	// Check token is valid and up to date
	token, err := application.VerifyToken(params.Token, r.App.Config.JWT.Refresh.Secret)
	if err != nil {
		return nil, err
	}
	// Extract userID claim in the token
	claims, err := application.ExtractTokenMetadata(token, []application.JwtClaimKey{application.UserIdClaim, application.UserAgentIdClaim})
	if err != nil {
		return nil, err
	}

	td, err := r.App.CreateTokens(claims[application.UserIdClaim], claims[application.UserAgentIdClaim])
	if err != nil {
		return nil, err
	}

	if claims[application.UserAgentIdClaim] != string(params.Agent.ID) {
		return nil, resolverErrUnauthorized(errors.New("Invalid session"))
	}

	u, err := r.App.Models.User.GetById(claims[application.UserIdClaim])
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return nil, resolverErrNotFound(err)
		default:
			return nil, resolverErrDatabaseOperation(err)
		}
	}

	if err = r.App.Models.User.InsertOrUpdateUserSession(
		&user.Session{
			ID:            string(params.Agent.ID),
			DeactivatedAt: time.Unix(td.RefreshExp, 0),
			IP:            params.Agent.IP,
			Name:          params.Agent.Name,
			Location:      params.Agent.Location,
			UserID:        u.ID,
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
	Agent AgentParams
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
