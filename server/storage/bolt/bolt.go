package bolt

import (
	"context"
	"fmt"

	"github.com/n101661/maney/pkg/utils"
	"github.com/n101661/maney/server/storage"
	bolt "go.etcd.io/bbolt"
)

const (
	userBucket  = "users"
	tokenBucket = "tokens"
)

type Storage struct {
	db *bolt.DB

	opts *options
}

func New(path string, opts ...utils.Option[options]) (storage.Storage, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open bolt db: %w", err)
	}

	s := &Storage{
		db:   db,
		opts: utils.ApplyOptions(defaultOptions(), opts),
	}

	if err := s.init(); err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	return s, nil
}

func (s *Storage) init() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(userBucket)); err != nil {
			return fmt.Errorf("failed to create user bucket: %w", err)
		}

		if _, err := tx.CreateBucketIfNotExists([]byte(tokenBucket)); err != nil {
			return fmt.Errorf("failed to create token bucket: %w", err)
		}

		return nil
	})
}

func (s *Storage) CreateUser(_ context.Context, user *storage.User) error {
	return create(s.db, userBucket, user.ID, user, s.opts)
}

func (s *Storage) GetUser(_ context.Context, userID string) (*storage.User, error) {
	return get[storage.User](s.db, userBucket, userID, s.opts)
}

func (s *Storage) CreateToken(_ context.Context, token *storage.Token) error {
	return create(s.db, tokenBucket, token.ID, token, s.opts)
}

func (s *Storage) GetToken(_ context.Context, tokenID string) (*storage.Token, error) {
	return get[storage.Token](s.db, tokenBucket, tokenID, s.opts)
}

func (s *Storage) DeleteToken(_ context.Context, tokenID string) (*storage.Token, error) {
	return delete[storage.Token](s.db, tokenBucket, tokenID, s.opts)
}

func (s *Storage) CreateConfig(ctx context.Context, config *storage.UserConfig) (*storage.UserConfig, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *Storage) GetConfig(ctx context.Context, id string) (*storage.UserConfig, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *Storage) UpdateConfig(ctx context.Context, config *storage.UserConfig) error {
	return fmt.Errorf("not implemented")
}

func (s *Storage) DeleteConfig(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func create(db *bolt.DB, bucketID, key string, value any, opts *options) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketID))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketID)
		}

		if bucket.Get([]byte(key)) != nil {
			return storage.ErrExists
		}

		data, err := opts.marshalValue(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}

		return bucket.Put([]byte(key), data)
	})
}

func get[T any](db *bolt.DB, bucketID, key string, opts *options) (*T, error) {
	var value T
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketID))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketID)
		}

		data := bucket.Get([]byte(key))
		if data == nil {
			return storage.ErrNotFound
		}

		return opts.unmarshalValue(data, &value)
	})
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func delete[T any](db *bolt.DB, bucketID, key string, opts *options) (*T, error) {
	var value T
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketID))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketID)
		}

		data := bucket.Get([]byte(key))
		if data == nil {
			return nil
		}

		if err := opts.unmarshalValue(data, &value); err != nil {
			return fmt.Errorf("failed to unmarshal value: %w", err)
		}

		return bucket.Delete([]byte(key))
	})
	if err != nil {
		return nil, err
	}
	return &value, nil
}
