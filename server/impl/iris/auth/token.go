package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	irisJWT "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

var (
	method = irisJWT.SigningMethodES512
)

type TokenClaims struct {
	UserID string
	Name   string
}

func newGenerateTokenFunc(secretKey string) GenerateTokenFunc {
	return func(claims TokenClaims, expiry time.Time) (token string, err error) {
		return irisJWT.NewTokenWithClaims(method, irisJWT.MapClaims(irisJWT.MapClaims{
			"user_id": claims.UserID,
			"name":    claims.Name,
			"exp":     jwt.NewNumericDate(expiry),
		})).SignedString(secretKey)
	}
}

func newValidateTokenFunc(secretKey string) ValidateTokenFunc {
	return irisJWT.New(irisJWT.Config{
		ValidationKeyGetter: func(*irisJWT.Token) (interface{}, error) {
			return secretKey, nil
		},
		SigningMethod: method,
		Expiration:    true,
	}).Serve
}

func getTokenClaims(ctx iris.Context) TokenClaims {
	token, ok := ctx.Values().Get(irisJWT.DefaultContextKey).(*irisJWT.Token)
	if !ok {
		return TokenClaims{}
	}

	claims, ok := token.Claims.(irisJWT.MapClaims)
	if !ok {
		return TokenClaims{}
	}
	return TokenClaims{
		UserID: claims["user_id"].(string),
		Name:   claims["name"].(string),
	}
}
