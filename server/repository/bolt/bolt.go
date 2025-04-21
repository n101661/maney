package bolt

import (
	"fmt"

	"go.etcd.io/bbolt"
	"golang.org/x/exp/constraints"

	"github.com/n101661/maney/server/repository"
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

func GetUserBucket(tx *bbolt.Tx, bucketName, userID string) (*bbolt.Bucket, error) {
	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return nil, fmt.Errorf("bucket %s not found", bucketName)
	}

	userBucket := bucket.Bucket([]byte(userID))
	return userBucket, nil
}

func GetUserBucketOrCreate(tx *bbolt.Tx, bucketName, userID string) (*bbolt.Bucket, error) {
	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return nil, fmt.Errorf("bucket %s not found", bucketName)
	}

	userBucket, err := bucket.CreateBucketIfNotExists([]byte(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to create and get %s account bucket: %v", userID, err)
	}
	return userBucket, nil
}

func NextSequence[T constraints.Signed](bucket *bbolt.Bucket) (T, error) {
	seq, err := bucket.NextSequence()
	if err != nil {
		return T(0), fmt.Errorf("failed to get next sequence of bucket: %w", err)
	}

	id := T(seq)
	if id < 0 || uint64(id) != seq {
		return T(0), fmt.Errorf("sequence overflow, origin: %d, to_T: %d", seq, id)
	}

	return id, nil
}
