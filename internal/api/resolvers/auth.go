package resolvers

import (
	"context"
	"time"

	"github.com/brice-74/golang-base-api/internal/api/application"
	"github.com/brice-74/golang-base-api/internal/domains/user"
	"github.com/brice-74/golang-base-api/pkg/validator"
	"github.com/graph-gophers/graphql-go"
	"github.com/ventu-io/go-shortid"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUserAccount: register a new user account
func (r Root) RegisterUserAccount(ctx context.Context, params RegisterUserAccountParams) (*UserAccountResolver, error) {
	uctx := r.App.UserFromContext(ctx)
	if !uctx.User.IsAnonymous() {
		return nil, resolverErrUnauthorized
	}

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
		if err == user.ErrDuplicateEmail {
			return nil, resolverError{Code: errDatabaseOperation, Message: err.Error(), StatusCode: 500}
		}
		return nil, err
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
func (r Root) LoginUserAccount(ctx context.Context, params LoginUserAccountParams) (*TokensUserAccountResolver, error) {
	uctx := r.App.UserFromContext(ctx)
	if !uctx.User.IsAnonymous() {
		return nil, resolverErrUnauthorized
	}

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
			return nil, resolverError{Code: errNotFound, Message: err.Error()}
		} else {
			return nil, err
		}
	}
	// check password
	if err = bcrypt.CompareHashAndPassword([]byte(uReg.Password), []byte(uEntry.Password)); err != nil {
		return nil, resolverError{Code: errInvalidCredentials, Message: "incorrect password"}
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
		return nil, err
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

func (r Root) RefreshUserAccount(ctx context.Context, params RefreshUserAccountParams) (*TokensUserAccountResolver, error) {
	uctx := r.App.UserFromContext(ctx)
	if uctx.Token.IsAccess {
		return nil, resolverErrUnauthorized
	}

	td, err := r.App.CreateTokens(uctx.User.ID, uctx.SessionID)
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
			UserID:        uctx.User.ID,
		},
	); err != nil {
		return nil, err
	}

	return &TokensUserAccountResolver{app: r.App, tokens: user.Tokens{
		Access:  td.AccessToken,
		Refresh: td.RefreshToken,
	}}, nil
}

type RefreshUserAccountParams struct {
	Agent AgentParams
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
