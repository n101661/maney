package users

import (
	"context"
	"errors"
	"time"
)

// Service errors.
var (
	ErrUserExists                    = errors.New("user exists")
	ErrUserNotFoundOrInvalidPassword = errors.New("user not found or invalid password")
	ErrInvalidToken                  = errors.New("invalid token")
	ErrTokenExpired                  = errors.New("token is expired")
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

	// ValidateRefreshToken validates if the refresh token is valid or not. It returns:
	//  - ErrInvalidToken if the refresh token is invalid
	//  - ErrTokenExpired if the refresh token is expired
	ValidateRefreshToken(ctx context.Context, r *ValidateRefreshTokenRequest) (*ValidateRefreshTokenReply, error)
}

type LoginRequest struct {
	UserID   string
	Password string
}

type LoginReply struct {
	AccessToken  *Token
	RefreshToken *Token
}

type TokenClaims struct {
	UserID string
	Nonce  int
}

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

type ValidateAccessTokenReply struct{}

type ValidateRefreshTokenRequest struct {
	TokenID string
}

type ValidateRefreshTokenReply struct{}
