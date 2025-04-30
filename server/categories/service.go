package categories

import (
	"context"
	"fmt"

	"github.com/n101661/maney/server/repository"
)

var (
	ErrDataInsufficient = fmt.Errorf("data insufficient")
	ErrCategoryNotFound = fmt.Errorf("category not found")
)

type Service interface {
	// Create returns ErrDataInsufficient if any of fields of CreateRequest is zero-value.
	Create(context.Context, *CreateRequest) (*CreateReply, error)
	// List returns ErrDataInsufficient if any of fields of ListRequest is zero-value,
	List(context.Context, *ListRequest) (*ListReply, error)
	// Update returns error:
	//  - ErrDataInsufficient if any of fields of UpdateRequest is zero-value,
	//  - ErrCategoryNotFound if the category does not exist.
	Update(context.Context, *UpdateRequest) (*UpdateReply, error)
	// Delete returns error:
	//  - ErrDataInsufficient if any of fields of UpdateRequest is zero-value,
	//  - ErrCategoryNotFound if the category does not exist.
	Delete(context.Context, *DeleteRequest) (*DeleteReply, error)
}

type Type = repository.CategoryType

type CreateRequest struct {
	UserID   string
	Type     Type
	Category *BaseCategory
}

type CreateReply struct {
	Type     Type
	Category *Category
}

type ListRequest struct {
	UserID string
	Type   Type
}

type ListReply struct {
	Categories []*Category
}

type UpdateRequest struct {
	UserID           string
	CategoryPublicID string

	Category *BaseCategory
}

type UpdateReply struct {
	Category *Category
}

type DeleteRequest struct {
	UserID           string
	CategoryPublicID string
}

type DeleteReply struct{}

type Category = repository.Category

type BaseCategory = repository.BaseCategory
