package user

import (
	"fmt"
	"strings"

	"github.com/brice-74/golang-base-api/pkg/validator"
)

func (user User) ValidateProfilNameEntry(v *validator.Validator) {
	profil_name := strings.TrimSpace(user.ProfilName)
	v.SetTmpKeyError("profil name")
	v.Check(profil_name != "", "must be provided")
	v.Check(len(profil_name) <= 32, "must have maximum of 32 characters")
}

func (user User) ValidatePasswordEntry(v *validator.Validator) {
	pass := strings.TrimSpace(user.Password)
	v.SetTmpKeyError("password")
	v.Check(pass != "", "must be provided")
	v.Check(len(pass) <= 255, "must have maximum of 255 characters")
	v.Check(len(pass) >= 8, "must have minimum of 8 characters")
	v.Check(validator.LowercaseRX(1, 255).MatchString(pass), "must have minimum of 1 lowercase")
	v.Check(validator.UppercaseRX(1, 255).MatchString(pass), "must have minimum of 1 uppercase")
	v.Check(validator.DigitRX(1, 255).MatchString(pass), "must have minimum of 1 number")
	v.Check(validator.SpecialCharRX(1, 255).MatchString(pass), "must have minimum of 1 character")
}

func (user User) ValidateEmailEntry(v *validator.Validator) {
	email := strings.TrimSpace(user.Email)
	v.SetTmpKeyError("email")
	v.Check(email != "", "must be provided")
	v.Check(len(email) <= 255, "must have maximum of 255 characters")
	v.Check(validator.EmailRX.MatchString(email), "must be a valid address")
}

func (user User) ValidateRolesEntry(v *validator.Validator, accepted_roles Roles) {
	v.SetTmpKeyError("roles")
	v.Check(len(user.Roles) == len(accepted_roles), fmt.Sprintf("must be provided: %s", strings.Join(accepted_roles.strings(), ", ")))
	for _, role := range user.Roles {
		v.Check(validator.In(string(role), accepted_roles.strings()...), fmt.Sprintf("invalid role : %s", role))
	}
}
