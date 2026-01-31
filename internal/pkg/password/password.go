package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyPass = errors.New("password is empty")
)

func Hash(pass string) (string, error) {
	if pass == "" {
		return "", ErrEmptyPass
	}

	bhash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	return string(bhash), err
}

func CompareHash(pass string, passHash string) bool {
	if pass == "" || passHash == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(pass))
	return err == nil
}
