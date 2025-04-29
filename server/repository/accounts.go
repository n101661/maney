package repository

import (
	"context"

	"github.com/shopspring/decimal"
)

type AccountRepository interface {
	// Create creates accounts of specific user and return error:
	//  - ErrDataExists if the data exists
	// or returns Account model with id.
	Create(context.Context, *CreateAccountsRequest) ([]*Account, error)
	// List returns accounts, it returns error:
	//  - ErrDataNotFound if there is no account satisfied filter conditions.
	List(context.Context, *ListAccountsRequest) (*ListAccountsReply, error)
	// Update updates non-zero value fields on specific account of the user, it returns error:
	//  - ErrDataNotFound if the account does not exist.
	Update(context.Context, *UpdateAccountRequest) (*Account, error)
	// Delete returns error:
	//  - ErrDataNotFound if the account does not exist.
	Delete(context.Context, *DeleteAccountsRequest) ([]*Account, error)
}

type CreateAccountsRequest struct {
	UserID   string
	Accounts []*BaseCreateAccount
}

type BaseCreateAccount struct {
	PublicID string
	*BaseAccount
}

type Account struct {
	ID       int32
	PublicID string
	*BaseAccount
	Balance decimal.Decimal
}

type BaseAccount struct {
	Name           string
	IconID         int32
	InitialBalance decimal.Decimal
}

type ListAccountsRequest struct {
	UserID          string
	AccountPublicID *string
}

type ListAccountsReply struct {
	Accounts []*Account
}

type UpdateAccountRequest struct {
	UserID          string
	AccountPublicID string

	Account      *BaseAccount
	BalanceDelta *decimal.Decimal
}

type DeleteAccountsRequest struct {
	AccountPublicIDs []string
	UserID           string
}
