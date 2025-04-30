package main

import (
	"fmt"
	"io"
	"time"

	_ "github.com/lib/pq"
	"xorm.io/xorm"
	"xorm.io/xorm/names"

	"github.com/n101661/maney/server/accounts"
	"github.com/n101661/maney/server/categories"
	"github.com/n101661/maney/server/repository"
	"github.com/n101661/maney/server/repository/postgres"
	"github.com/n101661/maney/server/shops"
	"github.com/n101661/maney/server/users"
)

type Repositories struct {
	User     repository.UserRepository
	Account  repository.AccountRepository
	Category repository.CategoryRepository
	Shop     repository.ShopRepository

	closer io.Closer
}

func newRepository(config *StorageConfig) (*Repositories, error) {
	if config.Postgres != nil {
		return newPostgresRepositories(config.Postgres)
	}
	return nil, fmt.Errorf("required storage setting")
}

func (repos *Repositories) Close() error {
	return repos.closer.Close()
}

func newPostgresRepositories(config *postgres.Config) (*Repositories, error) {
	connString := fmt.Sprintf(
		"postgresql://%s:%d/%s?user=%s&password=%s&sslmode=%s&",
		config.Host,
		config.Port,
		config.Database,
		config.User,
		config.Password,
		"disable",
	)

	engine, err := newXormEngine("postgres", connString, &xormEngineOptions{
		Schema:          config.Schema,
		ConnMaxIdleTime: config.ConnMaxIdleTime,
		ConnMaxLifetime: config.ConnMaxLifetime,
		MaxIdleConns:    config.MaxIdleConns,
		MaxOpenConns:    config.MaxOpenConns,
	})
	if err != nil {
		return nil, err
	}

	userRepo, err := users.NewPostgresRepository(engine)
	if err != nil {
		return nil, fmt.Errorf("failed to initial user repository: %v", err)
	}

	accountRepo, err := accounts.NewPostgresRepository(engine)
	if err != nil {
		return nil, fmt.Errorf("failed to initial account repository: %v", err)
	}

	categoryRepo, err := categories.NewPostgresRepository(engine)
	if err != nil {
		return nil, fmt.Errorf("failed to initial category repository: %v", err)
	}

	shopRepo, err := shops.NewPostgresRepository(engine)
	if err != nil {
		return nil, fmt.Errorf("failed to initial shop repository: %v", err)
	}

	return &Repositories{
		User:     userRepo,
		Account:  accountRepo,
		Category: categoryRepo,
		Shop:     shopRepo,
		closer:   engine,
	}, nil
}

type xormEngineOptions struct {
	Schema          string
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

func newXormEngine(driverName, dataSourceName string, opts *xormEngineOptions) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to create xorm engine: %v", err)
	}

	if opts.Schema != "" {
		engine.SetSchema(opts.Schema)
	}
	engine.SetColumnMapper(names.LintGonicMapper)
	engine.SetConnMaxIdleTime(opts.ConnMaxIdleTime)
	engine.SetConnMaxLifetime(opts.ConnMaxLifetime)
	engine.SetMaxIdleConns(opts.MaxIdleConns)
	engine.SetMaxOpenConns(opts.MaxOpenConns)

	err = engine.Sync(
		postgres.UsersModel{},
		postgres.TokensModel{},
		postgres.AccountsModel{},
		postgres.CategoriesModel{},
		postgres.ShopsModel{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sync tables: %v", err)
	}

	return engine, nil
}
