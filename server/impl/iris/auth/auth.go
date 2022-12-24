package auth

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type (
	GenerateTokenFunc    func(claims TokenClaims) (token string, err error)
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

func NewAuthentication(secretKey string) *Authentication {
	return &Authentication{
		GenerateToken:    newGenerateTokenFunc(secretKey),
		ValidateToken:    newValidateTokenFunc(secretKey),
		GetTokenClaims:   getTokenClaims,
		EncryptPassword:  encryptPassword,
		ValidatePassword: validatePassword,
	}
}
