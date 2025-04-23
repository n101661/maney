package accounts

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.etcd.io/bbolt"

	"github.com/n101661/maney/pkg/utils"
	"github.com/n101661/maney/pkg/utils/types"
	"github.com/n101661/maney/server/repository"
	"github.com/n101661/maney/server/repository/bolt"
)

const (
	accountBucket      = "account"
	publicIDToIDBucket = "publicIDToIDBucket"
)

type boltRepository struct {
	db *bbolt.DB

	opts *bolt.Options
}

func NewBoltRepository(path string, opts ...utils.Option[bolt.Options]) (repository.AccountRepository, error) {
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
		if _, err := tx.CreateBucketIfNotExists([]byte(accountBucket)); err != nil {
			return fmt.Errorf("failed to create account bucket: %w", err)
		}

		if _, err := tx.CreateBucketIfNotExists([]byte(publicIDToIDBucket)); err != nil {
			return fmt.Errorf("failed to create %s bucket: %w", publicIDToIDBucket, err)
		}

		return nil
	})
}

func (repo *boltRepository) Create(ctx context.Context, r *repository.CreateAccountsRequest) ([]*repository.Account, error) {
	rows := make([]*repository.Account, len(r.Accounts))
	err := repo.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := bolt.GetUserBucketOrCreate(tx, accountBucket, r.UserID)
		if err != nil {
			return err
		}

		publicIDBucket := tx.Bucket([]byte(publicIDToIDBucket))

		for i, account := range r.Accounts {
			publicID := []byte(account.PublicID)

			if id := publicIDBucket.Get(publicID); id != nil {
				return repository.ErrDataExists
			}

			id, err := bolt.NextSequence[int32](bucket)
			if err != nil {
				return err
			}

			bid := types.Int32ToBytes(id)

			if bucket.Get(bid) != nil {
				return fmt.Errorf("duplicated next sequence[%d]", id)
			}

			data, err := repo.opts.MarshalValue(&boltAccountModel{
				PublicID:    account.PublicID,
				BaseAccount: account.BaseAccount,
				Balance:     account.InitialBalance,
			})
			if err != nil {
				return fmt.Errorf("failed to marshal account: %w", err)
			}

			if err := publicIDBucket.Put(publicID, bid); err != nil {
				return fmt.Errorf("failed to put data to %s bucket: %w", publicIDToIDBucket, err)
			}
			if err := bucket.Put(bid, data); err != nil {
				return err
			}

			rows[i] = &repository.Account{
				ID:          id,
				PublicID:    account.PublicID,
				BaseAccount: lo.ToPtr(*account.BaseAccount),
				Balance:     account.InitialBalance,
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (repo *boltRepository) List(ctx context.Context, r *repository.ListAccountsRequest) (*repository.ListAccountsReply, error) {
	rows := []*repository.Account{}
	err := repo.db.View(func(tx *bbolt.Tx) error {
		bucket, err := bolt.GetUserBucket(tx, accountBucket, r.UserID)
		if err != nil {
			return err
		}
		if bucket == nil {
			return nil
		}

		if r.AccountPublicID != nil {
			publicID := []byte(*r.AccountPublicID)

			bid := tx.Bucket([]byte(publicIDToIDBucket)).Get(publicID)
			if bid == nil {
				return nil
			}
			data := bucket.Get(bid)

			id, err := types.BytesToInt32(bid)
			if err != nil {
				return err
			}

			rows, err = addRow(rows, id, data, repo.opts.UnmarshalValue)
			return err
		}

		return bucket.ForEach(func(k, v []byte) error {
			id, err := types.BytesToInt32(k)
			if err != nil {
				return fmt.Errorf("invalid id of account: %w", err)
			}

			rows, err = addRow(rows, id, v, repo.opts.UnmarshalValue)
			return err
		})
	})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, repository.ErrDataNotFound
	}
	return &repository.ListAccountsReply{
		Accounts: rows,
	}, nil
}

func addRow(
	rows []*repository.Account,
	id int32,
	value []byte,
	UnmarshalValue func([]byte, interface{}) error,
) ([]*repository.Account, error) {
	var account boltAccountModel
	if err := UnmarshalValue(value, &account); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data of account: %w", err)
	}

	rows = append(rows, &repository.Account{
		ID:          id,
		PublicID:    account.PublicID,
		BaseAccount: account.BaseAccount,
		Balance:     account.Balance,
	})

	return rows, nil
}

func (repo *boltRepository) Update(ctx context.Context, r *repository.UpdateAccountRequest) (*repository.Account, error) {
	var res *repository.Account
	err := repo.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := bolt.GetUserBucketOrCreate(tx, accountBucket, r.UserID)
		if err != nil {
			return err
		}

		bid := tx.Bucket([]byte(publicIDToIDBucket)).Get([]byte(r.AccountPublicID))
		if bid == nil {
			return repository.ErrDataNotFound
		}

		data := bucket.Get(bid)
		if data == nil {
			return fmt.Errorf("found in %s bucket but not found in %s bucket", publicIDToIDBucket, accountBucket)
		}

		var current boltAccountModel
		if err := repo.opts.UnmarshalValue(data, &current); err != nil {
			return fmt.Errorf("failed to unmarshal data of account: %w", err)
		}

		if r.Account != nil {
			current.BaseAccount = lo.ToPtr(*r.Account)
		}
		if r.BalanceDelta != nil {
			current.Balance = current.Balance.Add(*r.BalanceDelta)
		}

		data, err = repo.opts.MarshalValue(&current)
		if err != nil {
			return fmt.Errorf("failed to marshal account: %w", err)
		}

		if err := bucket.Put(bid, data); err != nil {
			return fmt.Errorf("failed to update account: %v", err)
		}

		id, err := types.BytesToInt32(bid)
		if err != nil {
			return err
		}

		res = &repository.Account{
			ID:          id,
			PublicID:    r.AccountPublicID,
			BaseAccount: current.BaseAccount,
			Balance:     current.Balance,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *boltRepository) Delete(ctx context.Context, r *repository.DeleteAccountsRequest) ([]*repository.Account, error) {
	rows := make([]*repository.Account, len(r.AccountPublicIDs))
	err := repo.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := bolt.GetUserBucketOrCreate(tx, accountBucket, r.UserID)
		if err != nil {
			return err
		}

		publicIDBucket := tx.Bucket([]byte(publicIDToIDBucket))

		for i, publicID := range r.AccountPublicIDs {
			bPublicID := []byte(publicID)

			bid := publicIDBucket.Get(bPublicID)
			if bid == nil {
				return repository.ErrDataNotFound
			}

			data := bucket.Get(bid)
			if data == nil {
				return fmt.Errorf("found in %s bucket but not found in %s bucket", publicIDToIDBucket, accountBucket)
			}

			if err := publicIDBucket.Delete(bPublicID); err != nil {
				return fmt.Errorf("failed to delete from %s bucket: %w", publicIDToIDBucket, err)
			}
			if err := bucket.Delete(bid); err != nil {
				return fmt.Errorf("failed to delete account: %w", err)
			}

			var account boltAccountModel
			if err := repo.opts.UnmarshalValue(data, &account); err != nil {
				return fmt.Errorf("failed to unmarshal data of account: %w", err)
			}

			id, err := types.BytesToInt32(bid)
			if err != nil {
				return err
			}

			rows[i] = &repository.Account{
				ID:          id,
				PublicID:    publicID,
				BaseAccount: account.BaseAccount,
				Balance:     account.Balance,
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

type boltAccountModel struct {
	PublicID string
	*repository.BaseAccount
	Balance decimal.Decimal
}
