package auth

import (
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type (
	GenerateTokenFunc    func(claims TokenClaims, expiry time.Time) (token string, err error)
	ValidateTokenFunc    = context.Handler
	GetTokenClaimsFunc   func(iris.Context) TokenClaims
	EncryptPasswordFunc  func(v string) ([]byte, error)
	ValidatePasswordFunc func(expected, actual []byte) error
)

type Authentication struct {
	GenerateToken    GenerateTokenFunc
	ValidateToken    ValidateTokenFunc
	GetTokenClaims   GetTokenClaimsFunc
	EncryptPassword  EncryptPasswordFunc
	ValidatePassword ValidatePasswordFunc
}

func NewAuthentication(secretKey []byte, opts ...authenticationOption) *Authentication {
	var o authenticationOptions
	for _, opt := range opts {
		opt(&o)
	}

	return &Authentication{
		GenerateToken:    newGenerateTokenFunc(secretKey),
		ValidateToken:    newValidateTokenFunc(secretKey),
		GetTokenClaims:   getTokenClaims,
		EncryptPassword:  newEncryptPassword(o.passwordSaltRound),
		ValidatePassword: validatePassword,
	}
}

type authenticationOption func(*authenticationOptions)

func WithPasswordSaltRound(round int) authenticationOption {
	return func(o *authenticationOptions) {
		o.passwordSaltRound = round
	}
}

type authenticationOptions struct {
	passwordSaltRound int
}
