package shops

import (
	"context"
	"fmt"
)

var (
	ErrDataInsufficient = fmt.Errorf("data insufficient")
	ErrShopNotFound     = fmt.Errorf("shop not found")
)

type Service interface {
	// Create returns ErrDataInsufficient if any of fields of CreateRequest is zero-value.
	Create(context.Context, *CreateRequest) (*CreateReply, error)
	// List returns ErrDataInsufficient if any of fields of ListRequest is zero-value,
	List(context.Context, *ListRequest) (*ListReply, error)
	// Update returns error:
	//  - ErrDataInsufficient if any of fields of UpdateRequest is zero-value,
	//  - ErrShopNotFound if the shop does not exist.
	Update(context.Context, *UpdateRequest) (*UpdateReply, error)
	// Delete returns error:
	//  - ErrDataInsufficient if any of fields of DeleteRequest is zero-value,
	//  - ErrShopNotFound if the shop does not exist.
	Delete(context.Context, *DeleteRequest) (*DeleteReply, error)
}

type BaseShop struct {
	Name    string
	Address string
}

type Shop struct {
	ID       int32
	PublicID string
	*BaseShop
}

type CreateRequest struct {
	UserID string
	Shop   *BaseShop
}

type CreateReply struct {
	Shop *Shop
}

type ListRequest struct {
	UserID string
}

type ListReply struct {
	Shops []*Shop
}

type UpdateRequest struct {
	UserID       string
	ShopPublicID string
	Shop         *BaseShop
}

type UpdateReply struct {
	Shop *Shop
}

type DeleteRequest struct {
	UserID       string
	ShopPublicID string
}

type DeleteReply struct{}
