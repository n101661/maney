package auth

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/n101661/maney/pkg/services/auth/storage"
)

type TokenClaims = storage.TokenClaims

type Token struct {
	ID          string
	Claims      *TokenClaims
	ExpireAfter time.Duration
}

type accessTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
