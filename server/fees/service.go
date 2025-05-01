package fees

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

var (
	ErrDataInsufficient = fmt.Errorf("data insufficient")
	ErrFeeNotFound      = fmt.Errorf("fee not found")
)

type Service interface {
	// Create returns ErrDataInsufficient if any of fields of CreateRequest is zero-value.
	Create(context.Context, *CreateRequest) (*CreateReply, error)
	// List returns ErrDataInsufficient if any of fields of ListRequest is zero-value,
	List(context.Context, *ListRequest) (*ListReply, error)
	// Update returns error:
	//  - ErrDataInsufficient if any of fields of UpdateRequest is zero-value,
	//  - ErrFeeNotFound if the fee does not exist.
	Update(context.Context, *UpdateRequest) (*UpdateReply, error)
	// Delete returns error:
	//  - ErrDataInsufficient if any of fields of DeleteRequest is zero-value,
	//  - ErrFeeNotFound if the fee does not exist.
	Delete(context.Context, *DeleteRequest) (*DeleteReply, error)
}

type BaseFee struct {
	Name string

	// Type determines fee type:
	//  - 0: Rate
	//  - 1: Fixed
	Type  int8
	Rate  *decimal.Decimal
	Fixed *decimal.Decimal
}

type Fee struct {
	ID       int32
	PublicID string
	*BaseFee
}

type CreateRequest struct {
	UserID string
	Fee    *BaseFee
}

type CreateReply struct {
	Fee *Fee
}

type ListRequest struct {
	UserID string
}

type ListReply struct {
	Fees []*Fee
}

type UpdateRequest struct {
	UserID      string
	FeePublicID string
	Fee         *BaseFee
}

type UpdateReply struct {
	Fee *Fee
}

type DeleteRequest struct {
	UserID      string
	FeePublicID string
}

type DeleteReply struct{}
