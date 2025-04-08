package auth

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/n101661/maney/pkg/utils"
)

type options struct {
	saltPasswordRound        int
	refreshTokenExpireAfter  time.Duration
	accessTokenSigningMethod jwt.SigningMethod
	accessTokenExpireAfter   time.Duration
}

func defaultOptions() *options {
	return &options{
		saltPasswordRound:        10,
		refreshTokenExpireAfter:  24 * time.Hour * 30,
		accessTokenSigningMethod: jwt.SigningMethodHS256,
		accessTokenExpireAfter:   10 * time.Minute,
	}
}

func WithSaltPasswordRound(round int) utils.Option[options] {
	return func(o *options) {
		o.saltPasswordRound = round
	}
}

func WithRefreshTokenExpireAfter(duration time.Duration) utils.Option[options] {
	return func(o *options) {
		o.refreshTokenExpireAfter = duration
	}
}

func WithAccessTokenSigningMethod(method jwt.SigningMethod) utils.Option[options] {
	return func(o *options) {
		o.accessTokenSigningMethod = method
	}
}

func WithAccessTokenExpireAfter(duration time.Duration) utils.Option[options] {
	return func(o *options) {
		o.accessTokenExpireAfter = duration
	}
}
