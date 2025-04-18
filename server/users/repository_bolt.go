package users

import (
	"context"
	"fmt"

	"go.etcd.io/bbolt"

	"github.com/n101661/maney/pkg/utils"
	"github.com/n101661/maney/server/repository"
	"github.com/n101661/maney/server/repository/bolt"
)

const (
	userBucket  = "users"
	tokenBucket = "tokens"
)

type boltRepository struct {
	db *bbolt.DB

	opts *bolt.Options
}

func NewBoltRepository(path string, opts ...utils.Option[bolt.Options]) (repository.UserRepository, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open bolt db: %w", err)
	}

	s := &boltRepository{
		db:   db,
		opts: utils.ApplyOptions(bolt.DefaultOptions(), opts),
	}

	if err := s.init(); err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	return s, nil
}

func (s *boltRepository) init() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(userBucket)); err != nil {
			return fmt.Errorf("failed to create user bucket: %w", err)
		}

		if _, err := tx.CreateBucketIfNotExists([]byte(tokenBucket)); err != nil {
			return fmt.Errorf("failed to create token bucket: %w", err)
		}

		return nil
	})
}

func (s *boltRepository) CreateUser(_ context.Context, user *repository.UserModel) error {
	return bolt.Create(s.db, userBucket, user.ID, user, s.opts)
}

func (s *boltRepository) GetUser(_ context.Context, userID string) (*repository.UserModel, error) {
	return bolt.Get[repository.UserModel](s.db, userBucket, userID, s.opts)
}

func (s *boltRepository) UpdateUser(_ context.Context, user *repository.UserModel) error {
	if user.Password == nil && user.Config == nil {
		return fmt.Errorf("no fields to update")
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", userBucket)
		}

		data := bucket.Get([]byte(user.ID))
		if data == nil {
			return repository.ErrDataNotFound
		}

		var current repository.UserModel
		if err := s.opts.UnmarshalValue(data, &current); err != nil {
			return fmt.Errorf("failed to unmarshal existing user: %w", err)
		}

		if user.Password != nil {
			current.Password = user.Password
		}
		if user.Config != nil {
			current.Config = user.Config
		}

		data, err := s.opts.MarshalValue(&current)
		if err != nil {
			return fmt.Errorf("failed to marshal updated user: %w", err)
		}

		return bucket.Put([]byte(user.ID), data)
	})
}

func (s *boltRepository) CreateToken(_ context.Context, token *repository.TokenModel) error {
	return bolt.Create(s.db, tokenBucket, token.ID, token, s.opts)
}

func (s *boltRepository) GetToken(_ context.Context, tokenID string) (*repository.TokenModel, error) {
	return bolt.Get[repository.TokenModel](s.db, tokenBucket, tokenID, s.opts)
}

func (s *boltRepository) DeleteToken(_ context.Context, tokenID string) (*repository.TokenModel, error) {
	return bolt.Delete[repository.TokenModel](s.db, tokenBucket, tokenID, s.opts)
}

func (s *boltRepository) Close() error {
	return s.db.Close()
}
