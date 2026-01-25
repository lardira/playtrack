package password

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(pass string) (string, error) {
	bhash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	return string(bhash), err
}

func CompareHash(pass string, passHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(pass))
	return err == nil
}
