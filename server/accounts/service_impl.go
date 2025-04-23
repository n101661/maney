package accounts

import (
	"context"
	"errors"
	"fmt"

	"github.com/n101661/maney/pkg/utils"
	"github.com/n101661/maney/pkg/utils/slugid"
	"github.com/n101661/maney/server/repository"
	"github.com/samber/lo"
)

type service struct {
	repository repository.AccountRepository

	opts *accountServiceOptions
}

func NewService(
	repository repository.AccountRepository,
	opts ...utils.Option[accountServiceOptions],
) (Service, error) {
	return &service{
		repository: repository,
		opts:       utils.ApplyOptions(defaultAccountServiceOptions(), opts),
	}, nil
}

func (s *service) Create(ctx context.Context, r *CreateRequest) (*CreateReply, error) {
	rows, err := s.repository.Create(ctx, &repository.CreateAccountsRequest{
		UserID: r.UserID,
		Accounts: []*repository.BaseCreateAccount{
			parseBaseCreateAccount(r.Account, s.opts.genPublicID),
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
		UserID:          r.UserID,
		AccountPublicID: lo.ToPtr(r.AccountPublicID),
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	row, err := s.repository.Update(ctx, &repository.UpdateAccountRequest{
		UserID:          r.UserID,
		AccountPublicID: r.AccountPublicID,
		Account:         parseBaseAccount(r.Account),
		BalanceDelta:    lo.ToPtr(r.Account.InitialBalance.Sub(origin.Accounts[0].InitialBalance)),
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
		AccountPublicIDs: []string{r.AccountPublicID},
		UserID:           r.UserID,
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
		ID:       v.ID,
		PublicID: v.PublicID,
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

func parseBaseCreateAccount(v *BaseAccount, genPublicID func() string) *repository.BaseCreateAccount {
	if v == nil {
		return nil
	}
	return &repository.BaseCreateAccount{
		PublicID:    genPublicID(),
		BaseAccount: parseBaseAccount(v),
	}
}

type accountServiceOptions struct {
	genPublicID func() string
}

func defaultAccountServiceOptions() *accountServiceOptions {
	return &accountServiceOptions{
		genPublicID: func() string {
			return slugid.New("act", 11)
		},
	}
}

func WithAccountServiceGenPublicID(f func() string) utils.Option[accountServiceOptions] {
	return func(o *accountServiceOptions) {
		o.genPublicID = f
	}
}
