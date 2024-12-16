package auth

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID string
}

type accessTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
