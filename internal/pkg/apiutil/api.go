package apiutil

import "github.com/golang-jwt/jwt/v5"

const (
	RoleAdmin  = "admin"
	RolePlayer = "player"
)

var (
	DefaultSigningMethod = jwt.SigningMethodHS256

	OperationSecurity = []map[string][]string{
		{"bearer": {"JWT"}},
	}
)
