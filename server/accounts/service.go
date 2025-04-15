package accounts

import (
	"context"

	"github.com/shopspring/decimal"
)

type Service interface {
	Create(context.Context, *CreateRequest) (*CreateReply, error)
	List(context.Context, *ListRequest) (*ListReply, error)
	Update(context.Context, *UpdateRequest) (*UpdateReply, error)
	Delete(context.Context, *DeleteRequest) (*DeleteReply, error)
}

type Account struct {
	ID             int64
	Name           string
	IconID         int64
	InitialBalance decimal.Decimal
}

type CreateRequest struct {
	UserID  string
	Account *Account
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
	UserID  string
	Account *Account
}

type UpdateReply struct{}

type DeleteRequest struct {
	UserID    string
	AccountID int64
}

type DeleteReply struct{}
