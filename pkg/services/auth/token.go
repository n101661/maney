package auth

import (
	"time"
)

type TokenClaims struct {
	UserID      string
	ExpiryAfter time.Duration
	Nonce       int
}
