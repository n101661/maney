package bolt

import (
	"fmt"

	"github.com/n101661/maney/server/repository"

	"go.etcd.io/bbolt"
)

func Create(db *bbolt.DB, bucketID, key string, value any, opts *Options) error {
	return db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketID))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketID)
		}

		if bucket.Get([]byte(key)) != nil {
			return repository.ErrDataExists
		}

		data, err := opts.MarshalValue(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}

		return bucket.Put([]byte(key), data)
	})
}

func Get[T any](db *bbolt.DB, bucketID, key string, opts *Options) (*T, error) {
	var value T
	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketID))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketID)
		}

		data := bucket.Get([]byte(key))
		if data == nil {
			return repository.ErrDataNotFound
		}

		return opts.UnmarshalValue(data, &value)
	})
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func Delete[T any](db *bbolt.DB, bucketID, key string, opts *Options) (*T, error) {
	var value T
	err := db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketID))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketID)
		}

		data := bucket.Get([]byte(key))
		if data == nil {
			return repository.ErrDataNotFound
		}

		if err := opts.UnmarshalValue(data, &value); err != nil {
			return fmt.Errorf("failed to unmarshal value: %w", err)
		}

		return bucket.Delete([]byte(key))
	})
	if err != nil {
		return nil, err
	}
	return &value, nil
}
