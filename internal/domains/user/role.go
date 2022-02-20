package user

import "fmt"

const (
	RoleAnonymous Role = "ROLE_ANONYMOUS"
	RoleUser      Role = "ROLE_USER"
)

type Role string
type Roles []Role

// Scan allows custom type to be Scanned by databases, by implementing the Scanner interface.
func (r *Role) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		*r = Role(src.([]byte))
	default:
		return fmt.Errorf("cannot scan type UserRole with type %s", v)
	}

	return nil
}
