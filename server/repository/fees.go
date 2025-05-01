package repository

import (
	"context"

	"github.com/shopspring/decimal"
)

type FeeRepository interface {
	// Create creates fees of specific user and return error:
	//  - ErrDataExists if the data exists
	// or returns Fee model with id.
	Create(context.Context, *CreateFeesRequest) ([]*Fee, error)
	// List returns fees, it returns error:
	//  - ErrDataNotFound if there is no fee satisfied filter conditions.
	List(context.Context, *ListFeesRequest) (*ListFeesReply, error)
	// Update updates non-zero value fields on specific fee of the user, it returns error:
	//  - ErrDataNotFound if the fee does not exist.
	Update(context.Context, *UpdateFeeRequest) (*Fee, error)
	// Delete returns error:
	//  - ErrDataNotFound if the fee does not exist.
	Delete(context.Context, *DeleteFeesRequest) ([]*Fee, error)
}

type CreateFeesRequest struct {
	UserID string
	Fees   []*BaseCreateFee
}

type BaseCreateFee struct {
	PublicID string
	*BaseFee
}

type Fee struct {
	ID       int32
	PublicID string
	*BaseFee
}

type BaseFee struct {
	Name string

	Type  int8
	Rate  *decimal.Decimal
	Fixed *decimal.Decimal
}

type ListFeesRequest struct {
	UserID string
}

type ListFeesReply struct {
	Fees []*Fee
}

type UpdateFeeRequest struct {
	UserID      string
	FeePublicID string

	Fee *BaseFee
}

type DeleteFeesRequest struct {
	FeePublicIDs []string
	UserID       string
}
