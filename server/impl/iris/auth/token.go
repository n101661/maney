package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	irisJWT "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

const (
	claimUserID   = "maney_user_id"
	claimUserName = "maney_name"
)

var (
	method = irisJWT.SigningMethodHS512
)

type TokenClaims struct {
	UserID string
	Name   string
}

func newGenerateTokenFunc(secretKey []byte) GenerateTokenFunc {
	return func(claims TokenClaims, expiry time.Time) (token string, err error) {
		return irisJWT.NewTokenWithClaims(method, irisJWT.MapClaims(irisJWT.MapClaims{
			claimUserID:   claims.UserID,
			claimUserName: claims.Name,
			"exp":         jwt.NewNumericDate(expiry),
		})).SignedString(secretKey)
	}
}

func newValidateTokenFunc(secretKey []byte) ValidateTokenFunc {
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
		UserID: claims[claimUserID].(string),
		Name:   claims[claimUserName].(string),
	}
}
