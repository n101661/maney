package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	// Create creates the given user. It returns ErrExists if the user already exists.
	CreateUser(ctx context.Context, user *User) error
	// GetUser returns the specified user. It returns ErrNotFound if the user does not exist.
	GetUser(ctx context.Context, userID string) (*User, error)

	// Create creates the given token. It returns ErrExists if the token already exists.
	CreateToken(ctx context.Context, token *Token) error
	// GetToken returns the specified token. It returns ErrNotFound if the token does not exist.
	GetToken(ctx context.Context, tokenID string) (*Token, error)
	// DeleteToken deletes the specified token. It returns the deleted token if successful.
	DeleteToken(ctx context.Context, tokenID string) (*Token, error)

	io.Closer
}

type User struct {
	ID       string
	Password []byte
}

type Token struct {
	ID         string
	Claim      *TokenClaims
	ExpiryTime time.Time
}

type TokenClaims struct {
	UserID string
	Nonce  int
}
