package accounts

import (
	"context"

	"github.com/shopspring/decimal"
)

type Repository interface {
	Create(context.Context, []*AccountModel) ([]*AccountModel, error)
	List(context.Context, *ListAccountModelRequest) (*ListAccountModelReply, error)
	Update(context.Context, *AccountModel) (*AccountModel, error)
	Delete(context.Context, *DeleteAccountModelRequest) error
}

type AccountModel struct {
	ID             int64
	Name           string
	IconID         int64
	InitialBalance decimal.Decimal
	UserID         string
}

type ListAccountModelRequest struct {
	UserID string
}

type ListAccountModelReply struct {
	Accounts []*AccountModel
}

type DeleteAccountModelRequest struct {
	AccountID int64
	UserID    string
}
