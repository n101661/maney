package categories

import (
	"context"
	"fmt"

	"github.com/n101661/maney/server/repository"
)

var (
	ErrCategoryNotFound = fmt.Errorf("category not found")
)

type Service interface {
	Create(context.Context, *CreateRequest) (*CreateReply, error)
	List(context.Context, *ListRequest) (*ListReply, error)
	// Update returns error:
	//  - ErrCategoryNotFound if the category does not exist.
	Update(context.Context, *UpdateRequest) (*UpdateReply, error)
	// Delete returns error:
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
	UserID     string
	CategoryID int32

	Category *BaseCategory
}

type UpdateReply struct {
	Category *Category
}

type DeleteRequest struct {
	UserID     string
	CategoryID int32
}

type DeleteReply struct{}

type Category = repository.Category

type BaseCategory = repository.BaseCategory
