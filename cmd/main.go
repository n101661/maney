package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/n101661/maney/pkg/logger"
	"github.com/n101661/maney/pkg/services/auth"
	"github.com/n101661/maney/pkg/services/auth/storage/bolt"
	"github.com/n101661/maney/server/impl/iris"
	"go.uber.org/zap"
)

const configPath = "config.toml"

func main() {
	config, err := LoadConfig(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err = CreateDefaultConfig(configPath); err != nil {
				log.Fatalf("failed to create %s: %v", configPath, err)
			}
			fmt.Printf("the %s has been created, please setup first", configPath)
			return
		}
		log.Fatalf("failed to load config: %v", err)
	}

	logger, err := logger.New(&logger.Config{})
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	authStorage, err := bolt.New(config.Auth.BoltDBPath)
	if err != nil {
		logger.Fatal("failed to initial the storage of the authentication service", zap.Error(err))
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
		logger.Fatal("failed to initial the authentication service", zap.Error(err))
	}

	s := iris.NewServer(iris.Config{}, authService)
	if err := s.ListenAndServe(fmt.Sprintf("%s:%d", config.App.Host, config.App.Port)); err != nil {
		logger.Fatal("failed to listen and serve", zap.Error(err))
	}
}
