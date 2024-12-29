package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Development bool          `toml:"development" comment:"Enable development mode."`
	Level       zapcore.Level `toml:"level" comment:"Log level. It can be one of the following: debug, info, warn, error, dpanic, panic, fatal. The default is info."`
}

func New(config *Config) (*zap.Logger, error) {
	cfg := newConfig(config.Development)

	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	cfg.Level.SetLevel(config.Level)

	return cfg.Build()
}

func newConfig(development bool) (cfg zap.Config) {
	defer func() {
		if !development {
			cfg.DisableCaller = true
		}
	}()

	if development {
		return zap.NewDevelopmentConfig()
	}
	return zap.NewProductionConfig()
}
