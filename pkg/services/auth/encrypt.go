package auth

import (
	"crypto/sha512"

	"golang.org/x/crypto/bcrypt"
)

func encryptPassword(pwd string, saltRound int) ([]byte, error) {
	encrypted := encrypt([]byte(pwd))
	return bcrypt.GenerateFromPassword(encrypted, saltRound)
}

func encrypt(val []byte) []byte {
	h := sha512.New()
	h.Write(val)
	return h.Sum(nil)
}
