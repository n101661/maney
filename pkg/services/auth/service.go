package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/n101661/maney/pkg/models"
	"github.com/n101661/maney/pkg/services/auth/storage"
	"github.com/n101661/maney/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	// CreateUser creates a new user with the given data. If the user already exists it returns
	// ErrUserExists error.
	CreateUser(ctx context.Context, user *models.User) error
	// ValidateUser validates the user and the password. If the user does not exist or the password
	// is invalid it returns ErrUserNotFoundOrInvalidPassword error.
	ValidateUser(ctx context.Context, id, password string) error

	GenerateRefreshToken(ctx context.Context, claim *TokenClaims) (tokenID string, err error)
	// ValidateRefreshToken validates if the refresh token is valid or not. It returns:
	//  - ErrInvalidToken if the refresh token is invalid
	//  - ErrTokenExpired if the refresh token is expired
	ValidateRefreshToken(ctx context.Context, tokenID string) error

	GenerateAccessToken(ctx context.Context, claim *TokenClaims) (tokenID string, err error)
	// ValidateAccessToken validates if the access token is valid or not. It returns:
	//  - ErrInvalidToken if the access token is invalid
	//  - ErrTokenExpired if the access token is expired
	ValidateAccessToken(ctx context.Context, tokenID string) error
}

type service struct {
	storage storage.Storage
	secret  []byte

	opts *options
}

func NewService(storage storage.Storage, secret []byte, opts ...utils.Option[options]) Service {
	return &service{
		storage: storage,
		secret:  secret,
		opts:    utils.ApplyOptions(defaultOptions(), opts),
	}
}

func (s *service) CreateUser(ctx context.Context, user *models.User) error {
	password, err := encryptPassword(user.Password, s.opts.saltPasswordRound)
	if err != nil {
		return err
	}

	err = s.storage.CreateUser(ctx, &storage.User{
		ID:       user.ID,
		Password: password,
		Email:    user.Email,
	})
	if err != nil {
		if errors.Is(err, storage.ErrExists) {
			return ErrUserExists
		}
		return err
	}
	return nil
}

func (s *service) ValidateUser(ctx context.Context, id, password string) error {
	user, err := s.storage.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrUserNotFoundOrInvalidPassword
		}
		return err
	}

	if err = validatePassword(user.Password, password); err != nil {
		return ErrUserNotFoundOrInvalidPassword
	}
	return nil
}

func (s *service) GenerateRefreshToken(ctx context.Context, claim *TokenClaims) (string, error) {
	token, err := generateRefreshToken(claim, s.secret)
	if err != nil {
		return "", err
	}

	err = s.storage.CreateToken(ctx, &storage.Token{
		ID:         token,
		ExpiryTime: time.Now().Add(s.opts.refreshTokenExpireAfter),
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *service) ValidateRefreshToken(ctx context.Context, tokenID string) error {
	token, err := s.storage.GetToken(ctx, tokenID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrInvalidToken
		}
		return err
	}

	if time.Now().Before(token.ExpiryTime) {
		return nil
	}
	return ErrTokenExpired
}

func (s *service) GenerateAccessToken(ctx context.Context, claim *TokenClaims) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (s *service) ValidateAccessToken(ctx context.Context, tokenID string) error {
	return fmt.Errorf("not implemented")
}

func validatePassword(expected []byte, actual string) error {
	return bcrypt.CompareHashAndPassword(expected, encrypt([]byte(actual)))
}

func generateRefreshToken(claim *TokenClaims, secret []byte) (string, error) {
	payload, err := json.Marshal(claim)
	if err != nil {
		return "", err
	}

	hash := hmac.New(sha256.New, secret)

	n, err := hash.Write(payload)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt claim: %w", err)
	}
	if n != len(payload) {
		return "", fmt.Errorf("failed to encrypt claim: truncated data")
	}

	return base64.StdEncoding.EncodeToString(payload) + "." + base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}
