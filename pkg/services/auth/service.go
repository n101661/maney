package auth

import (
	"context"

	"github.com/n101661/maney/pkg/models"
)

type Service interface {
	// CreateUser creates a new user with the given data. If the user already exists it returns
	// ErrUserExists error.
	CreateUser(ctx context.Context, user *models.User) error
	// ValidateUser validates the user and the password. If the user does not exist or the password
	// is invalid it returns ErrUserNotFoundOrInvalidPassword error.
	ValidateUser(ctx context.Context, id, password string) error

	GenerateRefreshToken(ctx context.Context, claim *TokenClaims) (token string, err error)
	// ValidateRefreshToken validates if the refresh token is valid or not. It returns:
	//  - ErrInvalidToken if the refresh token is invalid
	//  - ErrTokenExpired if the refresh token is expired
	ValidateRefreshToken(ctx context.Context, token string) error

	GenerateAccessToken(ctx context.Context, claim *TokenClaims) (token string, err error)
	// ValidateAccessToken validates if the access token is valid or not. It returns:
	//  - ErrInvalidToken if the access token is invalid
	//  - ErrTokenExpired if the access token is expired
	ValidateAccessToken(ctx context.Context, token string) error
}
