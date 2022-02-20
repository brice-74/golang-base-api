package user

import (
	"time"
)

var AnonymousUser = &User{Roles: Roles{RoleAnonymous}}

type User struct {
	ID            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeactivatedAt time.Time
	Email         string
	Password      string
	Roles         Roles
	ProfilName    string
	ShortId       string
}

// IsAnonymous checks if a user instance is anonymous.
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type Tokens struct {
	Access  string
	Refresh string
}
