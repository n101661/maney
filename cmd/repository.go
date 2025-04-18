package main

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/n101661/maney/server/accounts"
	"github.com/n101661/maney/server/repository"
	"github.com/n101661/maney/server/users"
)

type Repositories struct {
	User    repository.UserRepository
	Account repository.AccountRepository
}

func newBoltRepositories(dbDir string) (*Repositories, error) {
	userRepo, err := users.NewBoltRepository(filepath.Join(dbDir, "users.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to initial user repository: %v", err)
	}

	accountRepo, err := accounts.NewBoltRepository(filepath.Join(dbDir, "accounts.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to initial account repository: %v", err)
	}

	return &Repositories{
		User:    userRepo,
		Account: accountRepo,
	}, nil
}

func (repos *Repositories) Close() error {
	var errs []error
	if err := repos.User.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := repos.Account.Close(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
