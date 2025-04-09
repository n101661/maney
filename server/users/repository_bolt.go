package users

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/n101661/maney/pkg/utils"
	bolt "go.etcd.io/bbolt"
)

const (
	userBucket  = "users"
	tokenBucket = "tokens"
)

type boltRepository struct {
	db *bolt.DB

	opts *boltOptions
}

func NewBoltRepository(path string, opts ...utils.Option[boltOptions]) (Repository, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open bolt db: %w", err)
	}

	s := &boltRepository{
		db:   db,
		opts: utils.ApplyOptions(defaultBoltOptions(), opts),
	}

	if err := s.init(); err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	return s, nil
}

func (s *boltRepository) init() error {
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

func (s *boltRepository) CreateUser(_ context.Context, user *UserModel) error {
	return create(s.db, userBucket, user.ID, user, s.opts)
}

func (s *boltRepository) GetUser(_ context.Context, userID string) (*UserModel, error) {
	return get[UserModel](s.db, userBucket, userID, s.opts)
}

func (s *boltRepository) UpdateUser(_ context.Context, user *UserModel) error {
	if user.Password == nil && user.Config == nil {
		return fmt.Errorf("no fields to update")
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", userBucket)
		}

		data := bucket.Get([]byte(user.ID))
		if data == nil {
			return ErrDataNotFound
		}

		var current UserModel
		if err := s.opts.unmarshalValue(data, &current); err != nil {
			return fmt.Errorf("failed to unmarshal existing user: %w", err)
		}

		if user.Password != nil {
			current.Password = user.Password
		}
		if user.Config != nil {
			current.Config = user.Config
		}

		data, err := s.opts.marshalValue(&current)
		if err != nil {
			return fmt.Errorf("failed to marshal updated user: %w", err)
		}

		return bucket.Put([]byte(user.ID), data)
	})
}

func (s *boltRepository) CreateToken(_ context.Context, token *TokenModel) error {
	return create(s.db, tokenBucket, token.ID, token, s.opts)
}

func (s *boltRepository) GetToken(_ context.Context, tokenID string) (*TokenModel, error) {
	return get[TokenModel](s.db, tokenBucket, tokenID, s.opts)
}

func (s *boltRepository) DeleteToken(_ context.Context, tokenID string) (*TokenModel, error) {
	return delete[TokenModel](s.db, tokenBucket, tokenID, s.opts)
}

func (s *boltRepository) Close() error {
	return s.db.Close()
}

func create(db *bolt.DB, bucketID, key string, value any, opts *boltOptions) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketID))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketID)
		}

		if bucket.Get([]byte(key)) != nil {
			return ErrDataExists
		}

		data, err := opts.marshalValue(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}

		return bucket.Put([]byte(key), data)
	})
}

func get[T any](db *bolt.DB, bucketID, key string, opts *boltOptions) (*T, error) {
	var value T
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketID))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketID)
		}

		data := bucket.Get([]byte(key))
		if data == nil {
			return ErrDataNotFound
		}

		return opts.unmarshalValue(data, &value)
	})
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func delete[T any](db *bolt.DB, bucketID, key string, opts *boltOptions) (*T, error) {
	var value T
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketID))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketID)
		}

		data := bucket.Get([]byte(key))
		if data == nil {
			return ErrDataNotFound
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

type boltOptions struct {
	marshalValue   func(interface{}) ([]byte, error)
	unmarshalValue func([]byte, interface{}) error
}

func defaultBoltOptions() *boltOptions {
	return &boltOptions{
		marshalValue:   json.Marshal,
		unmarshalValue: json.Unmarshal,
	}
}

func WithMarshaller(marshaler func(interface{}) ([]byte, error)) func(*boltOptions) {
	return func(o *boltOptions) {
		o.marshalValue = marshaler
	}
}

func WithUnmarshaler(unmarshaler func([]byte, interface{}) error) func(*boltOptions) {
	return func(o *boltOptions) {
		o.unmarshalValue = unmarshaler
	}
}
