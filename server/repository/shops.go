package repository

import (
	"context"
)

type ShopRepository interface {
	// Create creates shops of specific user and return error:
	//  - ErrDataExists if the data exists
	// or returns Shop model with id.
	Create(context.Context, *CreateShopsRequest) ([]*Shop, error)
	// List returns shops, it returns error:
	//  - ErrDataNotFound if there is no shop satisfied filter conditions.
	List(context.Context, *ListShopsRequest) (*ListShopsReply, error)
	// Update updates non-zero value fields on specific shop of the user, it returns error:
	//  - ErrDataNotFound if the shop does not exist.
	Update(context.Context, *UpdateShopRequest) (*Shop, error)
	// Delete returns error:
	//  - ErrDataNotFound if the shop does not exist.
	Delete(context.Context, *DeleteShopsRequest) ([]*Shop, error)
}

type CreateShopsRequest struct {
	UserID string
	Shops  []*BaseCreateShop
}

type BaseCreateShop struct {
	PublicID string
	*BaseShop
}

type Shop struct {
	ID       int32
	PublicID string
	*BaseShop
}

type BaseShop struct {
	Name    string
	Address string
}

type ListShopsRequest struct {
	UserID string
}

type ListShopsReply struct {
	Shops []*Shop
}

type UpdateShopRequest struct {
	UserID       string
	ShopPublicID string

	Shop *BaseShop
}

type DeleteShopsRequest struct {
	ShopPublicIDs []string
	UserID        string
}
