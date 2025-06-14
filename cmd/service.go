package main

import (
	"fmt"
	"time"

	"github.com/n101661/maney/server/accounts"
	"github.com/n101661/maney/server/categories"
	"github.com/n101661/maney/server/fees"
	"github.com/n101661/maney/server/shops"
	"github.com/n101661/maney/server/users"
)

type Services struct {
	User     users.Service
	Account  accounts.Service
	Category categories.Service
	Shop     shops.Service
	Fee      fees.Service
}

func newServices(repos *Repositories, authConfig *AuthServiceConfig) (*Services, error) {
	user, err := users.NewService(
		repos.User,
		[]byte(authConfig.RefreshTokenSigningKey),
		[]byte(authConfig.AccessTokenSigningKey),
		users.WithRefreshTokenExpireAfter(time.Duration(authConfig.RefreshTokenExpireAfter)),
		users.WithAccessTokenExpireAfter(time.Duration(authConfig.AccessTokenExpireAfter)),
		users.WithSaltPasswordRound(authConfig.SaltPasswordRound),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initial the user service: %v", err)
	}

	account, err := accounts.NewService(repos.Account)
	if err != nil {
		return nil, fmt.Errorf("failed to initial the account service: %v", err)
	}

	category, err := categories.NewService(repos.Category)
	if err != nil {
		return nil, fmt.Errorf("failed to initial the category service: %v", err)
	}

	shop, err := shops.NewService(repos.Shop)
	if err != nil {
		return nil, fmt.Errorf("failed to initial the shop service: %v", err)
	}

	fee, err := fees.NewService(repos.Fee)
	if err != nil {
		return nil, fmt.Errorf("failed to initial the fee service: %v", err)
	}

	return &Services{
		User:     user,
		Account:  account,
		Category: category,
		Shop:     shop,
		Fee:      fee,
	}, nil
}
