package accounts

import (
	"context"
	"errors"
	"fmt"

	"github.com/n101661/maney/server/repository"
	"github.com/samber/lo"
)

type service struct {
	repository repository.AccountRepository
}

func NewService(
	repository repository.AccountRepository,
) (Service, error) {
	return &service{
		repository: repository,
	}, nil
}

func (s *service) Create(ctx context.Context, r *CreateRequest) (*CreateReply, error) {
	rows, err := s.repository.Create(ctx, &repository.CreateAccountsRequest{
		UserID: r.UserID,
		Accounts: []*repository.BaseAccount{
			parseBaseAccount(r.Account),
		},
	})
	if err != nil {
		return nil, err
	}

	return &CreateReply{
		Account: parseAccount(rows[0]),
	}, nil
}

func (s *service) List(ctx context.Context, r *ListRequest) (*ListReply, error) {
	reply, err := s.repository.List(ctx, &repository.ListAccountsRequest{
		UserID: r.UserID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return &ListReply{
				Accounts: []*Account{},
			}, nil
		}
		return nil, err
	}

	return &ListReply{
		Accounts: lo.Map(reply.Accounts, func(item *repository.Account, _ int) *Account {
			return parseAccount(item)
		}),
	}, nil
}

func (s *service) Update(ctx context.Context, r *UpdateRequest) (*UpdateReply, error) {
	if r.Account == nil {
		return nil, fmt.Errorf("nothing to update")
	}

	origin, err := s.repository.List(ctx, &repository.ListAccountsRequest{
		UserID:    r.UserID,
		AccountID: lo.ToPtr(r.AccountID),
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	row, err := s.repository.Update(ctx, &repository.UpdateAccountRequest{
		UserID:       r.UserID,
		AccountID:    r.AccountID,
		Account:      parseBaseAccount(r.Account),
		BalanceDelta: lo.ToPtr(r.Account.InitialBalance.Sub(origin.Accounts[0].InitialBalance)),
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	return &UpdateReply{
		Account: parseAccount(row),
	}, nil
}

func (s *service) Delete(ctx context.Context, r *DeleteRequest) (*DeleteReply, error) {
	_, err := s.repository.Delete(ctx, &repository.DeleteAccountsRequest{
		AccountIDs: []int32{r.AccountID},
		UserID:     r.UserID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	return &DeleteReply{}, nil
}

func parseAccount(v *repository.Account) *Account {
	return &Account{
		ID: v.ID,
		BaseAccount: &BaseAccount{
			Name:           v.Name,
			IconID:         v.IconID,
			InitialBalance: v.InitialBalance,
		},
		Balance: v.Balance,
	}
}

func parseBaseAccount(v *BaseAccount) *repository.BaseAccount {
	if v == nil {
		return nil
	}
	return &repository.BaseAccount{
		Name:           v.Name,
		IconID:         v.IconID,
		InitialBalance: v.InitialBalance,
	}
}
