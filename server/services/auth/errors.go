package auth

import (
	"errors"
)

var (
	ErrUserExists                    = errors.New("user exists")
	ErrUserNotFoundOrInvalidPassword = errors.New("user not found or invalid password")
	ErrInvalidToken                  = errors.New("invalid token")
	ErrTokenExpired                  = errors.New("token is expired")
)
