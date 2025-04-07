package storage

import (
	"context"
)

type Storage interface {
	CreateConfig(ctx context.Context, config *UserConfig) (*UserConfig, error)
	GetConfig(ctx context.Context, id string) (*UserConfig, error)
	UpdateConfig(ctx context.Context, config *UserConfig) error
	DeleteConfig(ctx context.Context, id string) error
}

type UserConfig struct {
	ID                          string
	CompareItemsInDifferentShop bool
	CompareItemsInSameShop      bool
}
