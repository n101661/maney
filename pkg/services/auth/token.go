package auth

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID string
}

type Token struct {
	ID          string
	Claims      *TokenClaims
	ExpireAfter time.Duration
}

type accessTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
