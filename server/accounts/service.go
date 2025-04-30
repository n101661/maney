package accounts

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

var (
	ErrDataInsufficient = fmt.Errorf("data insufficient")
	ErrAccountNotFound  = fmt.Errorf("account not found")
)

type Service interface {
	// Create returns ErrDataInsufficient if any of fields of CreateRequest is zero-value.
	Create(context.Context, *CreateRequest) (*CreateReply, error)
	// List returns ErrDataInsufficient if any of fields of ListRequest is zero-value,
	List(context.Context, *ListRequest) (*ListReply, error)
	// Update returns error:
	//  - ErrDataInsufficient if any of fields of UpdateRequest is zero-value,
	//  - ErrAccountNotFound if the account does not exist.
	Update(context.Context, *UpdateRequest) (*UpdateReply, error)
	// Delete returns error:
	//  - ErrDataInsufficient if any of fields of UpdateRequest is zero-value,
	//  - ErrAccountNotFound if the account does not exist.
	Delete(context.Context, *DeleteRequest) (*DeleteReply, error)
}

type BaseAccount struct {
	Name           string
	IconID         int32
	InitialBalance decimal.Decimal
}

type Account struct {
	ID       int32
	PublicID string
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
	UserID          string
	AccountPublicID string
	Account         *BaseAccount
}

type UpdateReply struct {
	Account *Account
}

type DeleteRequest struct {
	UserID          string
	AccountPublicID string
}

type DeleteReply struct{}
