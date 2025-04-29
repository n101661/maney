package accounts

import (
	"context"

	"github.com/n101661/maney/server/repository"
	"github.com/n101661/maney/server/repository/postgres"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"xorm.io/xorm"
)

type postgresRepository struct {
	engine *xorm.Engine
}

func NewPostgresRepository(engine *xorm.Engine) (repository.AccountRepository, error) {
	return &postgresRepository{
		engine: engine,
	}, nil
}

func (repo *postgresRepository) Create(ctx context.Context, r *repository.CreateAccountsRequest) ([]*repository.Account, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	rows := lo.Map(r.Accounts, func(item *repository.BaseCreateAccount, _ int) *postgres.AccountsModel {
		return &postgres.AccountsModel{
			PublicID: item.PublicID,
			UserID:   r.UserID,
			Data: &postgres.BaseAccount{
				BaseAccount: lo.ToPtr(*item.BaseAccount),
			},
			Balance: decimal.NewNullDecimal(item.BaseAccount.InitialBalance),
		}
	})
	_, err := session.Insert(rows)
	if err != nil {
		if postgres.UniqueViolationError(err) {
			return nil, repository.ErrDataExists
		}
		return nil, err
	}
	return lo.Map(rows, func(item *postgres.AccountsModel, _ int) *repository.Account {
		return toAccount(item)
	}), nil
}

func (repo *postgresRepository) List(ctx context.Context, r *repository.ListAccountsRequest) (*repository.ListAccountsReply, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	var rows []*postgres.AccountsModel
	err := session.Find(&rows, &postgres.AccountsModel{
		PublicID: lo.FromPtr(r.AccountPublicID),
		UserID:   r.UserID,
	})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, repository.ErrDataNotFound
	}
	return &repository.ListAccountsReply{
		Accounts: lo.Map(rows, func(item *postgres.AccountsModel, _ int) *repository.Account {
			return toAccount(item)
		}),
	}, nil
}

func (repo *postgresRepository) Update(ctx context.Context, r *repository.UpdateAccountRequest) (*repository.Account, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	row := postgres.AccountsModel{
		PublicID: r.AccountPublicID,
		UserID:   r.UserID,
	}
	has, err := session.Get(&row)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, repository.ErrDataNotFound
	}

	var bean postgres.AccountsModel
	if r.Account != nil {
		row.Data = &postgres.BaseAccount{
			BaseAccount: r.Account,
		}
		bean.Data = row.Data
	}
	if r.BalanceDelta != nil {
		row.Balance.Decimal = row.Balance.Decimal.Add(*r.BalanceDelta)
		bean.Balance = row.Balance
	}

	affected, err := session.Update(&bean, &postgres.AccountsModel{
		ID: row.ID,
	})
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, repository.ErrDataNotFound
	}

	return toAccount(&row), nil
}

func (repo *postgresRepository) Delete(ctx context.Context, r *repository.DeleteAccountsRequest) ([]*repository.Account, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	session.Where("user_id = ?", r.UserID)
	if len(r.AccountPublicIDs) > 0 {
		session.In("public_id", r.AccountPublicIDs)
	}

	var rows []*postgres.AccountsModel
	err := session.Find(&rows)
	if err != nil {
		return nil, err
	}

	if len(r.AccountPublicIDs) > 0 && len(rows) != len(r.AccountPublicIDs) {
		return nil, repository.ErrDataNotFound
	}
	if len(r.AccountPublicIDs) == 0 && len(rows) == 0 {
		return nil, repository.ErrDataNotFound
	}

	_, err = session.In("id", lo.Map(rows, func(item *postgres.AccountsModel, _ int) any {
		return item.ID
	})).Table(&postgres.AccountsModel{}).Delete()
	if err != nil {
		return nil, err
	}

	return lo.Map(rows, func(item *postgres.AccountsModel, _ int) *repository.Account {
		return toAccount(item)
	}), nil
}

func toAccount(item *postgres.AccountsModel) *repository.Account {
	return &repository.Account{
		ID:          item.ID,
		PublicID:    item.PublicID,
		BaseAccount: item.Data.BaseAccount,
		Balance:     item.Balance.Decimal,
	}
}
