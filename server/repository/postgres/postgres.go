package postgres

import (
	"time"
)

type Config struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Database string `toml:"database"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Schema   string `toml:"schema"`

	ConnMaxIdleTime time.Duration `toml:"conn-max-idle-time" comment:"The value <= 0 means connections are not closed."`
	ConnMaxLifetime time.Duration `toml:"conn-max-lifetime" comment:"The value <= 0 means connections are not closed."`
	MaxIdleConns    int           `toml:"max-idle-conns" comment:"The value <= 0 means no idle connections are retained."`
	MaxOpenConns    int           `toml:"max-open-conns" comment:"The value <= 0 means there is no limit on the number of open connections."`
}
