package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/n101661/maney/server/accounts"
	"github.com/n101661/maney/server/impl/iris"
	"github.com/n101661/maney/server/users"
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

	userRepo, err := users.NewBoltRepository(filepath.Join(config.Auth.BoltDBDir, "users.db"))
	if err != nil {
		fmt.Printf("failed to initial user repository: %v", err)
		os.Exit(1)
	}
	defer userRepo.Close()

	userService, err := users.NewService(
		userRepo,
		[]byte(config.Auth.RefreshTokenSigningKey),
		[]byte(config.Auth.AccessTokenSigningKey),
		users.WithRefreshTokenExpireAfter(time.Duration(config.Auth.RefreshTokenExpireAfter)),
		users.WithAccessTokenExpireAfter(time.Duration(config.Auth.AccessTokenExpireAfter)),
		users.WithSaltPasswordRound(config.Auth.SaltPasswordRound),
	)
	if err != nil {
		fmt.Printf("failed to initial the user service: %v", err)
		os.Exit(1)
	}
	userController := users.NewIrisController(userService)

	accountRepo, err := accounts.NewBoltRepository(filepath.Join(config.Auth.BoltDBDir, "accounts.db"))
	if err != nil {
		fmt.Printf("failed to initial account repository: %v", err)
		os.Exit(1)
	}
	defer userRepo.Close()

	accountService, err := accounts.NewService(accountRepo)
	if err != nil {
		fmt.Printf("failed to initial the account service: %v", err)
		os.Exit(1)
	}
	accountController := accounts.NewIrisController(accountService)

	s := iris.NewServer(config.App.Config, &iris.Controllers{
		User:    userController,
		Account: accountController,
	})
	if err := s.ListenAndServe(fmt.Sprintf("%s:%d", config.App.Host, config.App.Port)); err != nil {
		fmt.Printf("failed to listen and serve: %v", err)
		os.Exit(1)
	}
}
