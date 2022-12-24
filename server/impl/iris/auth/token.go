package auth

import (
	"time"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

var (
	method = jwt.SigningMethodES512
)

type TokenClaims struct {
	UserID string
	Name   string
	Expiry time.Time
}

func newGenerateTokenFunc(secretKey string) GenerateTokenFunc {
	return func(claims TokenClaims) (token string, err error) {
		return jwt.NewTokenWithClaims(method, jwt.MapClaims(jwt.MapClaims{
			"user_id": claims.UserID,
			"name":    claims.Name,
			"expiry":  claims.Expiry,
		})).SignedString(secretKey)
	}
}

func newValidateTokenFunc(secretKey string) ValidateTokenFunc {
	return jwt.New(jwt.Config{
		ValidationKeyGetter: func(*jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
		SigningMethod: method,
		Expiration:    true,
	}).Serve
}

func getTokenClaims(ctx iris.Context) TokenClaims {
	token, ok := ctx.Values().Get(jwt.DefaultContextKey).(*jwt.Token)
	if !ok {
		return TokenClaims{}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return TokenClaims{}
	}
	return TokenClaims{
		UserID: claims["user_id"].(string),
		Name:   claims["name"].(string),
		Expiry: claims["expiry"].(time.Time),
	}
}
