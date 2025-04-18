package accounts

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

var (
	ErrAccountNotFound = fmt.Errorf("account not found")
)

type Service interface {
	Create(context.Context, *CreateRequest) (*CreateReply, error)
	List(context.Context, *ListRequest) (*ListReply, error)
	// Update returns error:
	//  - ErrAccountNotFound if the account does not exist.
	Update(context.Context, *UpdateRequest) (*UpdateReply, error)
	// Delete returns error:
	//  - ErrAccountNotFound if the account does not exist.
	Delete(context.Context, *DeleteRequest) (*DeleteReply, error)
}

type BaseAccount struct {
	Name           string
	IconID         int32
	InitialBalance decimal.Decimal
}

type Account struct {
	ID int32
	*BaseAccount
	// Balance is the current balance.
	Balance decimal.Decimal
}

type CreateRequest struct {
	UserID  string
	Account *BaseAccount
}

type CreateReply struct {
	Account *Account
}

type ListRequest struct {
	UserID string
}

type ListReply struct {
	Accounts []*Account
}

type UpdateRequest struct {
	UserID    string
	AccountID int32
	Account   *BaseAccount
}

type UpdateReply struct {
	Account *Account
}

type DeleteRequest struct {
	UserID    string
	AccountID int32
}

type DeleteReply struct{}
