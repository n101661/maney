package accounts

import (
	"context"
	"fmt"

	"go.etcd.io/bbolt"
	"golang.org/x/exp/constraints"

	"github.com/n101661/maney/pkg/utils"
	"github.com/n101661/maney/server/internal/repository"
	"github.com/n101661/maney/server/internal/repository/bolt"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

const (
	accountBucket        = "account"
	accountBalanceBucket = "accountBalance"
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

		if _, err := tx.CreateBucketIfNotExists([]byte(accountBalanceBucket)); err != nil {
			return fmt.Errorf("failed to create account balance bucket: %w", err)
		}

		return nil
	})
}

func (repo *boltRepository) Create(ctx context.Context, r *repository.CreateAccountsRequest) ([]*repository.Account, error) {
	rows := make([]*repository.Account, len(r.Accounts))
	err := repo.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := getUserBucket(tx, r.UserID)
		if err != nil {
			return err
		}

		for i, account := range r.Accounts {
			id, err := nextSequenceAs[int32](bucket)
			if err != nil {
				return err
			}

			bid := int32ToBytes(id)

			if bucket.Get(bid) != nil {
				return repository.ErrDataExists
			}

			data, err := repo.opts.MarshalValue(&boltAccountModel{
				BaseAccount: account,
				Balance:     account.InitialBalance,
			})
			if err != nil {
				return fmt.Errorf("failed to marshal account: %w", err)
			}

			if err := bucket.Put(bid, data); err != nil {
				return err
			}

			rows[i] = &repository.Account{
				ID:          id,
				BaseAccount: lo.ToPtr(*account),
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
		bucket, err := getUserBucket(tx, r.UserID)
		if err != nil {
			return err
		}

		return bucket.ForEach(func(k, v []byte) error {
			id, err := bytesToInt32(k)
			if err != nil {
				return fmt.Errorf("invalid id of account: %w", err)
			}

			var account boltAccountModel
			if err := repo.opts.UnmarshalValue(v, &account); err != nil {
				return fmt.Errorf("failed to unmarshal data of account: %w", err)
			}

			rows = append(rows, &repository.Account{
				ID:          id,
				BaseAccount: account.BaseAccount,
				Balance:     account.Balance,
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
	return &repository.ListAccountsReply{
		Accounts: rows,
	}, nil
}

func (repo *boltRepository) Update(ctx context.Context, r *repository.UpdateAccountRequest) (*repository.Account, error) {
	var res *repository.Account
	err := repo.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := getUserBucket(tx, r.UserID)
		if err != nil {
			return err
		}

		bid := int32ToBytes(r.AccountID)

		data := bucket.Get(bid)
		if data == nil {
			return repository.ErrDataNotFound
		}

		var current boltAccountModel
		if err := repo.opts.UnmarshalValue(data, &current); err != nil {
			return fmt.Errorf("failed to unmarshal data of account: %w", err)
		}

		if r.Account != nil {
			current.BaseAccount = lo.ToPtr(*r.Account)
		}
		if r.BalanceDelta != nil {
			current.Balance.Add(*r.BalanceDelta)
		}

		data, err = repo.opts.MarshalValue(&current)
		if err != nil {
			return fmt.Errorf("failed to marshal account: %w", err)
		}

		if err := bucket.Put(bid, data); err != nil {
			return fmt.Errorf("failed to update account: %v", err)
		}

		res = &repository.Account{
			ID:          r.AccountID,
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
	rows := make([]*repository.Account, len(r.AccountIDs))
	err := repo.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := getUserBucket(tx, r.UserID)
		if err != nil {
			return err
		}

		for i, accountID := range r.AccountIDs {
			bid := int32ToBytes(accountID)

			data := bucket.Get(bid)
			if data == nil {
				return repository.ErrDataNotFound
			}

			if err := bucket.Delete(bid); err != nil {
				return fmt.Errorf("failed to delete account: %w", err)
			}

			var account boltAccountModel
			if err := repo.opts.UnmarshalValue(data, &account); err != nil {
				return fmt.Errorf("failed to unmarshal data of account: %w", err)
			}

			rows[i] = &repository.Account{
				ID:          accountID,
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

type boltAccountModel struct {
	*repository.BaseAccount
	Balance decimal.Decimal
}

func getUserBucket(tx *bbolt.Tx, userID string) (*bbolt.Bucket, error) {
	bucket := tx.Bucket([]byte(accountBucket))
	if bucket == nil {
		return nil, fmt.Errorf("bucket %s not found", accountBucket)
	}

	userBucket, err := bucket.CreateBucketIfNotExists([]byte(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to create and get %s account bucket", userID)
	}
	return userBucket, nil
}

func nextSequenceAs[T constraints.Signed](bucket *bbolt.Bucket) (T, error) {
	seq, err := bucket.NextSequence()
	if err != nil {
		return T(0), fmt.Errorf("failed to get next sequence of account bucket: %w", err)
	}

	id := T(seq)
	if id < 0 || uint64(id) != seq {
		return T(0), fmt.Errorf("sequence overflow, origin: %d, to_int32: %d", seq, id)
	}

	return id, nil
}

func int32ToBytes(v int32) []byte {
	res := make([]byte, 4)
	for i := range 4 {
		res[3-i] = byte(v)
		v >>= 8
	}
	return res
}

func bytesToInt32(v []byte) (int32, error) {
	if len(v) > 4 {
		return 0, fmt.Errorf("invalid int32: %v", v)
	}
	res := int32(0)
	for i, b := range v {
		res |= int32(b) << ((3 - i) * 8)
	}
	return res, nil
}
