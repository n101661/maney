package repository

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/samber/lo"
)

type CategoryRepository interface {
	// Create creates categories of specific user and return error:
	//  - ErrDataExists if the data exists
	// or returns Category model with id.
	Create(context.Context, *CreateCategoriesRequest) ([]*Category, error)
	// List returns categories, it returns error:
	//  - ErrDataNotFound if there is no account satisfied filter conditions.
	List(context.Context, *ListCategoriesRequest) (*ListCategoriesReply, error)
	// Update updates non-zero value fields on specific account of the user, it returns error:
	//  - ErrDataNotFound if the account does not exist.
	Update(context.Context, *UpdateCategoryRequest) (*Category, error)
	// Delete returns error:
	//  - ErrDataNotFound if the account does not exist.
	Delete(context.Context, *DeleteCategoriesRequest) ([]*Category, error)

	io.Closer
}

type CreateCategoriesRequest struct {
	UserID     string
	Type       CategoryType
	Categories []*BaseCreateCategory
}

type BaseCreateCategory struct {
	PublicID string
	*BaseCategory
}

type ListCategoriesRequest struct {
	UserID string
	Type   CategoryType
}

type ListCategoriesReply struct {
	Categories []*Category
}

type UpdateCategoryRequest struct {
	UserID           string
	CategoryPublicID string

	Category *BaseCategory
}

type DeleteCategoriesRequest struct {
	UserID            string
	CategoryPublicIDs []string
}

const (
	CategoryTypeExpense CategoryType = iota
	CategoryTypeIncome
)

type CategoryType uint8

func (t CategoryType) String() string {
	if s, ok := typeDescriptions[t]; ok {
		return s
	}
	return strconv.Itoa(int(t))
}

var typeDescriptions = map[CategoryType]string{
	CategoryTypeExpense: "expense",
	CategoryTypeIncome:  "income",
}

var descriptionsToType = lo.Invert(typeDescriptions)

func ToCategoryType(s string) (CategoryType, error) {
	if t, ok := descriptionsToType[s]; ok {
		return t, nil
	}
	return 0, fmt.Errorf("unknown type of category[%s]", s)
}

type Category struct {
	ID       int32
	PublicID string
	*BaseCategory
}

type BaseCategory struct {
	Name   string
	IconID int32
}
