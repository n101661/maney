package users

import (
	"context"
	"errors"
	"time"

	"github.com/n101661/maney/server/models"
	"github.com/n101661/maney/server/repository"
)

// Service errors.
var (
	ErrUserExists                    = errors.New("user exists")
	ErrUserNotFoundOrInvalidPassword = errors.New("user not found or invalid password")
	ErrInvalidToken                  = errors.New("invalid token")
	ErrTokenExpired                  = errors.New("token is expired")
	ErrResourceNotFound              = errors.New("resource not found")
)

type Service interface {
	// Login validates the user and the password. If the user does not exist or the password
	// is invalid it returns ErrUserNotFoundOrInvalidPassword error.
	// If the user is valid it returns a LoginResponse with the access and refresh tokens.
	Login(ctx context.Context, r *LoginRequest) (*LoginReply, error)

	// Logout revokes the token. It returns:
	//  - ErrInvalidToken if the token is invalid
	//  - ErrTokenExpired if the token is expired
	Logout(ctx context.Context, r *LogoutRequest) (*LogoutReply, error)

	// SignUp creates a new user with the given data. If the user already exists it returns
	// ErrUserExists error.
	SignUp(ctx context.Context, r *SignUpRequest) (*SignUpReply, error)

	// ValidateAccessToken validates if the access token is valid or not. It returns:
	//  - ErrInvalidToken if the access token is invalid
	//  - ErrTokenExpired if the access token is expired
	ValidateAccessToken(ctx context.Context, r *ValidateAccessTokenRequest) (*ValidateAccessTokenReply, error)

	// RefreshAccessToken validates if the refresh token is valid or not. It returns error:
	//  - ErrInvalidToken if the refresh token is invalid
	//  - ErrTokenExpired if the refresh token is expired
	// and returns newer access token.
	RefreshAccessToken(ctx context.Context, r *RefreshAccessTokenRequest) (*RefreshAccessTokenReply, error)

	// UpdateConfig updates the config, it returns:
	//  - ErrResourceNotFound if the user is not found
	UpdateConfig(context.Context, *UpdateConfigRequest) (*UpdateConfigReply, error)

	// GetConfig gets the config of the user, it returns:
	//  - ErrResourceNotFound if the user is not found
	GetConfig(context.Context, *GetConfigRequest) (*GetConfigReply, error)
}

type LoginRequest struct {
	UserID   string
	Password string
}

type LoginReply struct {
	AccessToken  *Token
	RefreshToken *Token
}

type TokenClaims = repository.TokenClaims

type Token struct {
	ID          string
	Claims      *TokenClaims
	ExpireAfter time.Duration
}

type LogoutRequest struct {
	RefreshTokenID string
}

type LogoutReply struct{}

type SignUpRequest struct {
	UserID   string
	Password string
}

type SignUpReply struct{}

type ValidateAccessTokenRequest struct {
	TokenID string
}

type ValidateAccessTokenReply struct {
	UserID string
}

type RefreshAccessTokenRequest struct {
	TokenID string
}

type RefreshAccessTokenReply struct {
	AccessToken  *Token
	RefreshToken *Token
}

type UpdateConfigRequest struct {
	UserID string
	Config *models.UserConfig
}

type UpdateConfigReply struct{}

type UserConfig = models.UserConfig

type GetConfigRequest struct {
	UserID string
}

type GetConfigReply struct {
	Data *UserConfig
}
