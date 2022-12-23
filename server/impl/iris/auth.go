package iris

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

type authentication struct {
	generateToken func(claims map[string]interface{}) (token string, err error)
	validateToken func(iris.Context)
}

func newAuthentication(secretKey string) *authentication {
	method := jwt.SigningMethodES512

	return &authentication{
		generateToken: func(claims map[string]interface{}) (token string, err error) {
			return jwt.NewTokenWithClaims(method, jwt.MapClaims(claims)).SignedString(secretKey)
		},
		validateToken: jwt.New(jwt.Config{
			ValidationKeyGetter: func(*jwt.Token) (interface{}, error) {
				return secretKey, nil
			},
			SigningMethod: method,
			Expiration:    true,
		}).Serve,
	}
}
