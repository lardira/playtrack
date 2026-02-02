package apiutil

import "github.com/golang-jwt/jwt/v5"

type Role string

const (
	RoleAdmin  Role = "admin"
	RolePlayer Role = "player"
)

func (r Role) String() string {
	return string(r)
}

var (
	DefaultSigningMethod = jwt.SigningMethodHS256

	OperationSecurity = []map[string][]string{
		{"bearer": {"JWT"}},
	}
)
