package iris

import (
	"crypto/sha512"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
)

type authentication struct {
	generateToken    func(claims map[string]interface{}) (token string, err error)
	validateToken    func(iris.Context)
	encryptPassword  func(v string) ([]byte, error)
	validatePassword func(expected, actual []byte) error
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
		encryptPassword: func(v string) ([]byte, error) {
			pwd := encryptPassword([]byte(v))
			return bcrypt.GenerateFromPassword(pwd, 8)
		},
		validatePassword: func(expected []byte, actual []byte) error {
			return bcrypt.CompareHashAndPassword(expected, encryptPassword(actual))
		},
	}
}
