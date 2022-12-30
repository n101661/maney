package auth

import (
	"crypto/sha512"

	"golang.org/x/crypto/bcrypt"
)

func newEncryptPassword(saltRound int) EncryptPasswordFunc {
	return func(val string) ([]byte, error) {
		pwd := encrypt([]byte(val))
		return bcrypt.GenerateFromPassword(pwd, saltRound)
	}
}

func validatePassword(expected []byte, actual []byte) error {
	return bcrypt.CompareHashAndPassword(expected, encrypt(actual))
}

func encrypt(val []byte) []byte {
	h := sha512.New()
	h.Write(val)
	return h.Sum(nil)
}
