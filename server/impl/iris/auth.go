package iris

import (
	"crypto/sha512"
	"time"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
)

type authentication struct {
	generateToken    func(claims tokenClaims) (token string, err error)
	validateToken    func(iris.Context)
	getTokenClaims   func(iris.Context) tokenClaims
	encryptPassword  func(v string) ([]byte, error)
	validatePassword func(expected, actual []byte) error
}

type tokenClaims struct {
	UserID string
	Name   string
	Expiry time.Time
}

func newAuthentication(secretKey string) *authentication {
	var (
		method          = jwt.SigningMethodES512
		encryptPassword = func(password []byte) []byte {
			h := sha512.New()
			h.Write(password)
			return h.Sum(nil)
		}
	)

	return &authentication{
		generateToken: func(claims tokenClaims) (token string, err error) {
			return jwt.NewTokenWithClaims(method, jwt.MapClaims(jwt.MapClaims{
				"user_id": claims.UserID,
				"name":    claims.Name,
				"expiry":  claims.Expiry,
			})).SignedString(secretKey)
		},
		validateToken: jwt.New(jwt.Config{
			ValidationKeyGetter: func(*jwt.Token) (interface{}, error) {
				return secretKey, nil
			},
			SigningMethod: method,
			Expiration:    true,
		}).Serve,
		getTokenClaims: func(ctx iris.Context) tokenClaims {
			token, ok := ctx.Values().Get(jwt.DefaultContextKey).(*jwt.Token)
			if !ok {
				return tokenClaims{}
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return tokenClaims{}
			}
			return tokenClaims{
				UserID: claims["user_id"].(string),
				Name:   claims["name"].(string),
				Expiry: claims["expiry"].(time.Time),
			}
		},
		encryptPassword: func(v string) ([]byte, error) {
			pwd := encryptPassword([]byte(v))
			return bcrypt.GenerateFromPassword(pwd, 8)
		},
		validatePassword: func(expected []byte, actual []byte) error {
			return bcrypt.CompareHashAndPassword(expected, encryptPassword(actual))
		},
	}
}
