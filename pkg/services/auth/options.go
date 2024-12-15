package auth

import (
	"time"

	"github.com/n101661/maney/pkg/utils"
)

type options struct {
	saltPasswordRound       int
	refreshTokenExpireAfter time.Duration
}

func defaultOptions() *options {
	return &options{
		saltPasswordRound:       10,
		refreshTokenExpireAfter: 24 * time.Hour * 30,
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
