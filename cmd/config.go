package main

import (
	"fmt"
	"os"
	"time"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/n101661/maney/pkg/encoding"
	"github.com/n101661/maney/server/impl/iris"
	"github.com/n101661/maney/server/impl/iris/config"
	"github.com/n101661/maney/server/repository/postgres"
)

type Config struct {
	App     *AppConfig         `toml:"application"`
	Auth    *AuthServiceConfig `toml:"authentication-service"`
	Storage *StorageConfig     `toml:"storage" comment:"Choose one of storage config as prefer storage. If you provide multiple settings, the system uses them in priority order: 'storage.postgres'."`
}

type AppConfig struct {
	Host string `toml:"host" comment:"Host of the application."`
	Port int    `toml:"port" comment:"Port of the application."`

	*iris.Config
}

type AuthServiceConfig struct {
	SaltPasswordRound       int               `toml:"salt-password-round" comment:"Number of rounds to salt the password. If the value is not provided or less than 0, the default is 10."`
	RefreshTokenSigningKey  string            `toml:"refresh-token-signing-key" comment:"Private key to sign the refresh token."`
	RefreshTokenExpireAfter encoding.Duration `toml:"refresh-token-expire-after" comment:"Period of the refresh token expiration. If the value is not provided, the default is 30 days."`
	AccessTokenSigningKey   string            `toml:"access-token-signing-key" comment:"Private key to sign the access token."`
	AccessTokenExpireAfter  encoding.Duration `toml:"access-token-expire-after" comment:"Period of the access token expiration. If the value is not provided, the default is 10 minutes."`
}

type StorageConfig struct {
	Postgres *postgres.Config `toml:"postgres" comment:"Connection settings of postgres."`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err = toml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func CreateDefaultConfig(path string) (err error) {
	config := &Config{
		App: &AppConfig{
			Host: "localhost",
			Port: 8080,
			Config: &iris.Config{
				LogLevel:    config.LogLevelInfo,
				CorsOrigins: []string{"*"},
			},
		},
		Auth: &AuthServiceConfig{
			SaltPasswordRound:       10,
			RefreshTokenSigningKey:  "THIS_IS_UNSECURE_SIGNED_KEY",
			RefreshTokenExpireAfter: encoding.Duration(24 * time.Hour * 30),
			AccessTokenSigningKey:   "THIS_IS_UNSECURE_SIGNED_KEY",
			AccessTokenExpireAfter:  encoding.Duration(10 * time.Minute),
		},
		Storage: &StorageConfig{
			Postgres: &postgres.Config{
				Host:            "",
				Port:            0,
				Database:        "",
				User:            "",
				Password:        "",
				Schema:          "",
				ConnMaxIdleTime: 30 * time.Minute,
				ConnMaxLifetime: time.Hour,
				MaxIdleConns:    2,
				MaxOpenConns:    2,
			},
		},
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create %s file: %w", path, err)
	}
	defer f.Close()
	defer func() {
		e := f.Sync()
		if err == nil {
			err = e
		}
	}()

	if err = toml.NewEncoder(f).Encode(config); err != nil {
		return fmt.Errorf("failed to write content to %s file:%w", path, err)
	}

	return nil
}
