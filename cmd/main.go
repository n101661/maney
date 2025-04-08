package main

import (
	"fmt"
	"os"
	"time"

	"github.com/n101661/maney/server/impl/iris"
	"github.com/n101661/maney/server/services/auth"
	"github.com/n101661/maney/server/storage/bolt"
)

const configPath = "config.toml"

func main() {
	config, err := LoadConfig(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err = CreateDefaultConfig(configPath); err != nil {
				fmt.Printf("failed to create %s: %v", configPath, err)
				os.Exit(1)
			}
			fmt.Printf("the %s has been created, please setup first", configPath)
			return
		}
		fmt.Printf("failed to load config: %v", err)
		os.Exit(1)
	}

	authStorage, err := bolt.New(config.Auth.BoltDBPath)
	if err != nil {
		fmt.Printf("failed to initial the storage of the authentication service: %v", err)
		os.Exit(1)
	}
	defer authStorage.Close()

	authService, err := auth.NewService(
		authStorage,
		[]byte(config.Auth.RefreshTokenSigningKey),
		[]byte(config.Auth.AccessTokenSigningKey),
		auth.WithRefreshTokenExpireAfter(time.Duration(config.Auth.RefreshTokenExpireAfter)),
		auth.WithAccessTokenExpireAfter(time.Duration(config.Auth.AccessTokenExpireAfter)),
		auth.WithSaltPasswordRound(config.Auth.SaltPasswordRound),
	)
	if err != nil {
		fmt.Printf("failed to initial the authentication service: %v", err)
		os.Exit(1)
	}

	s := iris.NewServer(config.App.Config, authService)
	if err := s.ListenAndServe(fmt.Sprintf("%s:%d", config.App.Host, config.App.Port)); err != nil {
		fmt.Printf("failed to listen and serve: %v", err)
		os.Exit(1)
	}
}
