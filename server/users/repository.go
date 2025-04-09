package users

import (
	"context"
	"errors"
	"io"
	"time"
)

// Repository errors.
var (
	ErrDataExists   = errors.New("the data exists")
	ErrDataNotFound = errors.New("the data is not found")
)

type Repository interface {
	// Create creates the given user. It returns ErrDataExists if the user already exists.
	CreateUser(ctx context.Context, user *UserModel) error
	// GetUser returns the specified user. It returns ErrDataNotFound if the user does not exist.
	GetUser(ctx context.Context, userID string) (*UserModel, error)
	// UpdateUser updates non-zero-value fields for specific user.
	// It returns ErrDataNotFound if the user does not exist.
	UpdateUser(ctx context.Context, user *UserModel) error

	// Create creates the given token. It returns ErrDataExists if the token already exists.
	CreateToken(ctx context.Context, token *TokenModel) error
	// GetToken returns the specified token. It returns ErrDataNotFound if the token does not exist.
	GetToken(ctx context.Context, tokenID string) (*TokenModel, error)
	// DeleteToken deletes the specified token. It returns the deleted token if successful.
	DeleteToken(ctx context.Context, tokenID string) (*TokenModel, error)

	io.Closer
}

type UserModel struct {
	ID       string
	Password []byte
	Config   *UserConfig
}

type UserConfig struct {
	CompareItemsInDifferentShop bool
	CompareItemsInSameShop      bool
}

type TokenModel struct {
	ID         string
	Claim      *TokenClaims
	ExpiryTime time.Time
}
