package categories

import (
	"context"
	"fmt"

	"go.etcd.io/bbolt"

	"github.com/n101661/maney/pkg/utils"
	"github.com/n101661/maney/pkg/utils/types"
	"github.com/n101661/maney/server/repository"
	"github.com/n101661/maney/server/repository/bolt"
	"github.com/samber/lo"
)

const (
	categoryBucket     = "category"
	publicIDToIDBucket = "publicIDToIDBucket"
)

type boltRepository struct {
	db *bbolt.DB

	opts *bolt.Options
}

func NewBoltRepository(path string, opts ...utils.Option[bolt.Options]) (repository.CategoryRepository, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open bolt db: %w", err)
	}

	repo := &boltRepository{
		db:   db,
		opts: utils.ApplyOptions(bolt.DefaultOptions(), opts),
	}

	if err := repo.init(); err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	return repo, nil
}

func (s *boltRepository) init() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(categoryBucket)); err != nil {
			return fmt.Errorf("failed to create category bucket: %w", err)
		}

		if _, err := tx.CreateBucketIfNotExists([]byte(publicIDToIDBucket)); err != nil {
			return fmt.Errorf("failed to create %s bucket: %w", publicIDToIDBucket, err)
		}

		return nil
	})
}

func (repo *boltRepository) Create(ctx context.Context, r *repository.CreateCategoriesRequest) ([]*Category, error) {
	rows := make([]*Category, len(r.Categories))
	err := repo.db.Update(func(tx *bbolt.Tx) error {
		userBucket, err := bolt.GetUserBucketOrCreate(tx, categoryBucket, r.UserID)
		if err != nil {
			return err
		}

		bucket, err := userBucket.CreateBucketIfNotExists([]byte{byte(r.Type)})
		if err != nil {
			return err
		}

		publicIDBucket := tx.Bucket([]byte(publicIDToIDBucket))

		for i, category := range r.Categories {
			publicID := []byte(category.PublicID)

			if id := publicIDBucket.Get(publicID); id != nil {
				return repository.ErrDataExists
			}

			id, err := bolt.NextSequence[int32](userBucket)
			if err != nil {
				return err
			}

			bid := types.Int32ToBytes(id)

			if userBucket.Get(bid) != nil {
				return fmt.Errorf("duplicated next sequence[%d]", id)
			}
			if bucket.Get(bid) != nil {
				return fmt.Errorf("not found in category bucket but found in category type bucket")
			}

			data, err := repo.opts.MarshalValue(&boltCategoryModel{
				PublicID: category.PublicID,
				BaseCategory: &repository.BaseCategory{
					Name:   category.Name,
					IconID: category.IconID,
				},
			})
			if err != nil {
				return fmt.Errorf("failed to marshal category: %w", err)
			}

			if err := publicIDBucket.Put(publicID, bid); err != nil {
				return err
			}
			if err := userBucket.Put(bid, []byte{byte(r.Type)}); err != nil {
				return err
			}
			if err := bucket.Put(bid, data); err != nil {
				return err
			}

			rows[i] = &repository.Category{
				ID:       id,
				PublicID: category.PublicID,
				BaseCategory: &repository.BaseCategory{
					Name:   category.Name,
					IconID: category.IconID,
				},
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (repo *boltRepository) List(ctx context.Context, r *repository.ListCategoriesRequest) (*repository.ListCategoriesReply, error) {
	rows := []*repository.Category{}
	err := repo.db.View(func(tx *bbolt.Tx) error {
		userBucket, err := bolt.GetUserBucket(tx, categoryBucket, r.UserID)
		if err != nil {
			return err
		}
		if userBucket == nil {
			return nil
		}

		bucket := userBucket.Bucket([]byte{byte(r.Type)})
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			id, err := types.BytesToInt32(k)
			if err != nil {
				return fmt.Errorf("invalid id of category: %w", err)
			}

			var category boltCategoryModel
			if err := repo.opts.UnmarshalValue(v, &category); err != nil {
				return fmt.Errorf("failed to unmarshal data of category: %w", err)
			}

			rows = append(rows, &repository.Category{
				ID:           id,
				PublicID:     category.PublicID,
				BaseCategory: category.BaseCategory,
			})

			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, repository.ErrDataNotFound
	}
	return &repository.ListCategoriesReply{
		Categories: rows,
	}, nil
}

func (repo *boltRepository) Update(ctx context.Context, r *repository.UpdateCategoryRequest) (*Category, error) {
	var res *repository.Category
	err := repo.db.Update(func(tx *bbolt.Tx) error {
		userBucket, err := bolt.GetUserBucketOrCreate(tx, categoryBucket, r.UserID)
		if err != nil {
			return err
		}

		bid := tx.Bucket([]byte(publicIDToIDBucket)).Get([]byte(r.CategoryPublicID))
		if bid == nil {
			return repository.ErrDataNotFound
		}

		bucketName := userBucket.Get(bid)
		if bucketName == nil {
			return fmt.Errorf("found in %s bucket but not found in category type bucket", publicIDToIDBucket)
		}

		bucket := userBucket.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("found in category bucket but there is no category type bucket[%v]", bucketName)
		}

		data := bucket.Get(bid)
		if data == nil {
			return fmt.Errorf("found in category bucket but not found in category type bucket")
		}

		var current boltCategoryModel
		if err := repo.opts.UnmarshalValue(data, &current); err != nil {
			return fmt.Errorf("failed to unmarshal data of category: %w", err)
		}

		if r.Category != nil {
			current.BaseCategory = lo.ToPtr(*r.Category)
		}

		data, err = repo.opts.MarshalValue(&current)
		if err != nil {
			return fmt.Errorf("failed to marshal category: %w", err)
		}

		if err := bucket.Put(bid, data); err != nil {
			return fmt.Errorf("failed to update category: %v", err)
		}

		id, err := types.BytesToInt32(bid)
		if err != nil {
			return err
		}

		res = &repository.Category{
			ID:           id,
			PublicID:     current.PublicID,
			BaseCategory: current.BaseCategory,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *boltRepository) Delete(ctx context.Context, r *repository.DeleteCategoriesRequest) ([]*Category, error) {
	rows := make([]*repository.Category, len(r.CategoryPublicIDs))
	err := repo.db.Update(func(tx *bbolt.Tx) error {
		userBucket, err := bolt.GetUserBucketOrCreate(tx, categoryBucket, r.UserID)
		if err != nil {
			return err
		}

		publicIDBucket := tx.Bucket([]byte(publicIDToIDBucket))

		for i, publicID := range r.CategoryPublicIDs {
			bPublicID := []byte(publicID)

			bid := publicIDBucket.Get(bPublicID)
			if bid == nil {
				return repository.ErrDataNotFound
			}

			bucketName := userBucket.Get(bid)
			if bucketName == nil {
				return fmt.Errorf("found in %s bucket but not found in %s bucket", publicIDToIDBucket, categoryBucket)
			}

			bucket := userBucket.Bucket(bucketName)
			if bucket == nil {
				return fmt.Errorf("found in category bucket but there is no category type bucket[%v]", bucketName)
			}

			data := bucket.Get(bid)
			if data == nil {
				return fmt.Errorf("found in category bucket but not found in category type bucket")
			}

			if err := publicIDBucket.Delete(bPublicID); err != nil {
				return fmt.Errorf("failed to delete from %s bucket: %w", publicIDToIDBucket, err)
			}
			if err := userBucket.Delete(bid); err != nil {
				return fmt.Errorf("failed to delete category from category bucket: %w", err)
			}
			if err := bucket.Delete(bid); err != nil {
				return fmt.Errorf("failed to delete category from category type bucket: %w", err)
			}

			var category boltCategoryModel
			if err := repo.opts.UnmarshalValue(data, &category); err != nil {
				return fmt.Errorf("failed to unmarshal data of category: %w", err)
			}

			id, err := types.BytesToInt32(bid)
			if err != nil {
				return err
			}

			rows[i] = &repository.Category{
				ID:           id,
				PublicID:     publicID,
				BaseCategory: category.BaseCategory,
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (repo *boltRepository) Close() error {
	return repo.db.Close()
}

type boltCategoryModel struct {
	PublicID string
	*repository.BaseCategory
}
