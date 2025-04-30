package users

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/n101661/maney/pkg/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/n101661/maney/server/repository"
)

type service struct {
	repository             repository.UserRepository
	accessTokenSigningKey  []byte
	refreshTokenSigningKey []byte

	opts *serviceOptions
}

func NewService(
	storage repository.UserRepository,
	accessTokenSigningKey []byte,
	refreshTokenSigningKey []byte,
	opts ...utils.Option[serviceOptions],
) (Service, error) {
	if len(accessTokenSigningKey) == 0 {
		return nil, errors.New("required access token signing key")
	}
	if len(refreshTokenSigningKey) == 0 {
		return nil, errors.New("required refresh token signing key")
	}

	return &service{
		repository:             storage,
		accessTokenSigningKey:  accessTokenSigningKey,
		refreshTokenSigningKey: refreshTokenSigningKey,
		opts:                   utils.ApplyOptions(defaultOptions(), opts),
	}, nil
}

func (s *service) Login(ctx context.Context, r *LoginRequest) (*LoginReply, error) {
	err := s.validateUser(ctx, r.UserID, r.Password)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.generateAccessToken(&TokenClaims{
		UserID: r.UserID,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(ctx, &TokenClaims{
		UserID: r.UserID,
	})
	if err != nil {
		return nil, err
	}

	return &LoginReply{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) validateUser(ctx context.Context, id, password string) error {
	user, err := s.repository.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return ErrUserNotFoundOrInvalidPassword
		}
		return err
	}

	if err = validatePassword(user.Password, password); err != nil {
		return err
	}
	return nil
}

func validatePassword(expected []byte, actual string) error {
	err := bcrypt.CompareHashAndPassword(expected, hashValue([]byte(actual)))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrUserNotFoundOrInvalidPassword
	}
	return err
}

type accessTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *service) generateAccessToken(claim *TokenClaims) (*Token, error) {
	token := jwt.NewWithClaims(s.opts.accessTokenSigningMethod, accessTokenClaims{
		UserID: claim.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.opts.accessTokenExpireAfter)),
		},
	})

	id, err := token.SignedString(s.accessTokenSigningKey)
	if err != nil {
		return nil, err
	}

	return &Token{
		ID:          id,
		Claims:      claim,
		ExpireAfter: s.opts.accessTokenExpireAfter,
	}, nil
}

func (s *service) generateRefreshToken(ctx context.Context, claim *TokenClaims) (*Token, error) {
	tokenID, err := generateRefreshToken(claim, s.refreshTokenSigningKey)
	if err != nil {
		return nil, err
	}

	err = s.repository.CreateToken(ctx, &repository.TokenModel{
		ID:         hashRefreshToken(tokenID),
		Claim:      claim,
		ExpiryTime: time.Now().Add(s.opts.refreshTokenExpireAfter),
	})
	if err != nil {
		return nil, err
	}
	return &Token{
		ID:          tokenID,
		Claims:      claim,
		ExpireAfter: s.opts.refreshTokenExpireAfter,
	}, nil
}

type refreshTokenClaims struct {
	*TokenClaims
	Seed uint64
}

func generateRefreshToken(claim *TokenClaims, signingKey []byte) (string, error) {
	payload, err := json.Marshal(refreshTokenClaims{
		TokenClaims: claim,
		Seed:        rand.Uint64(),
	})
	if err != nil {
		return "", err
	}

	hash := hmac.New(sha256.New, signingKey)

	n, err := hash.Write(payload)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt claim: %w", err)
	}
	if n != len(payload) {
		return "", fmt.Errorf("failed to encrypt claim: truncated data")
	}

	return base64.StdEncoding.EncodeToString(payload) + "." + base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}

func (s *service) Logout(ctx context.Context, r *LogoutRequest) (*LogoutReply, error) {
	_, err := s.revokeRefreshToken(ctx, hashRefreshToken(r.RefreshTokenID))
	if err != nil {
		return nil, err
	}
	return &LogoutReply{}, nil
}

func (s *service) revokeRefreshToken(ctx context.Context, tokenID string) (*repository.TokenModel, error) {
	token, err := s.repository.GetToken(ctx, tokenID)
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	now := time.Now()
	if now.After(token.ExpiryTime) {
		return nil, ErrTokenExpired
	}
	if token.RevokedAt != nil && now.After(*token.RevokedAt) {
		return nil, ErrInvalidToken
	}

	if err := s.repository.RevokeToken(ctx, tokenID); err != nil {
		return nil, err
	}
	return token, nil
}

func (s *service) SignUp(ctx context.Context, r *SignUpRequest) (*SignUpReply, error) {
	encryptedPassword, err := encryptPassword(r.Password, s.opts.saltPasswordRound)
	if err != nil {
		return nil, err
	}

	err = s.repository.CreateUser(ctx, &repository.UserModel{
		ID:       r.UserID,
		Password: encryptedPassword,
		Config:   &UserConfig{},
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataExists) {
			return nil, ErrUserExists
		}
		return nil, err
	}
	return &SignUpReply{}, nil
}

func encryptPassword(pwd string, saltRound int) ([]byte, error) {
	encrypted := hashValue([]byte(pwd))
	return bcrypt.GenerateFromPassword(encrypted, saltRound)
}

func (s *service) ValidateAccessToken(ctx context.Context, r *ValidateAccessTokenRequest) (*ValidateAccessTokenReply, error) {
	claims := accessTokenClaims{}

	_, err := jwt.ParseWithClaims(r.TokenID, &claims, func(t *jwt.Token) (interface{}, error) {
		return s.accessTokenSigningKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	return &ValidateAccessTokenReply{
		UserID: claims.UserID,
	}, nil
}

func (s *service) RefreshAccessToken(ctx context.Context, r *RefreshAccessTokenRequest) (*RefreshAccessTokenReply, error) {
	token, err := s.revokeRefreshToken(ctx, hashRefreshToken(r.TokenID))
	if err != nil {
		return nil, err
	}

	accessToken, err := s.generateAccessToken(&TokenClaims{
		UserID: token.Claim.UserID,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(ctx, &TokenClaims{
		UserID: token.Claim.UserID,
	})
	if err != nil {
		return nil, err
	}

	return &RefreshAccessTokenReply{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) UpdateConfig(ctx context.Context, r *UpdateConfigRequest) (*UpdateConfigReply, error) {
	err := s.repository.UpdateUser(ctx, &repository.UserModel{
		ID:     r.UserID,
		Config: r.Config,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}
	return &UpdateConfigReply{}, nil
}

func (s *service) GetConfig(ctx context.Context, r *GetConfigRequest) (*GetConfigReply, error) {
	user, err := s.repository.GetUser(ctx, r.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}
	return &GetConfigReply{
		Data: user.Config,
	}, nil
}

type serviceOptions struct {
	saltPasswordRound        int
	refreshTokenExpireAfter  time.Duration
	accessTokenSigningMethod jwt.SigningMethod
	accessTokenExpireAfter   time.Duration
	getNonce                 func() int
}

func defaultOptions() *serviceOptions {
	return &serviceOptions{
		saltPasswordRound:        10,
		refreshTokenExpireAfter:  24 * time.Hour * 30,
		accessTokenSigningMethod: jwt.SigningMethodHS256,
		accessTokenExpireAfter:   10 * time.Minute,
		getNonce: func() int {
			return int(time.Now().UnixNano()) % 9999
		},
	}
}

func WithSaltPasswordRound(round int) utils.Option[serviceOptions] {
	return func(o *serviceOptions) {
		o.saltPasswordRound = round
	}
}

func WithRefreshTokenExpireAfter(duration time.Duration) utils.Option[serviceOptions] {
	return func(o *serviceOptions) {
		o.refreshTokenExpireAfter = duration
	}
}

func WithAccessTokenSigningMethod(method jwt.SigningMethod) utils.Option[serviceOptions] {
	return func(o *serviceOptions) {
		o.accessTokenSigningMethod = method
	}
}

func WithAccessTokenExpireAfter(duration time.Duration) utils.Option[serviceOptions] {
	return func(o *serviceOptions) {
		o.accessTokenExpireAfter = duration
	}
}

func WithNonceGenerator(f func() int) utils.Option[serviceOptions] {
	return func(o *serviceOptions) {
		o.getNonce = f
	}
}

func hashValue(val []byte) []byte {
	h := sha512.New()
	h.Write(val)
	return h.Sum(nil)
}

func hashRefreshToken(token string) string {
	return base64.StdEncoding.EncodeToString(hashValue([]byte(token)))
}
