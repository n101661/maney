package categories

import (
	"context"
	"errors"
	"fmt"

	"github.com/n101661/maney/server/repository"
)

type service struct {
	repository repository.CategoryRepository
}

func NewService(
	repository repository.CategoryRepository,
) (Service, error) {
	return &service{
		repository: repository,
	}, nil
}

func (s *service) Create(ctx context.Context, r *CreateRequest) (*CreateReply, error) {
	if r.Category == nil {
		return nil, fmt.Errorf("nothing to create")
	}

	rows, err := s.repository.Create(ctx, &repository.CreateCategoriesRequest{
		UserID: r.UserID,
		Type:   r.Type,
		Categories: []*repository.BaseCategory{
			r.Category,
		},
	})
	if err != nil {
		return nil, err
	}

	return &CreateReply{
		Type:     r.Type,
		Category: rows[0],
	}, nil
}

func (s *service) List(ctx context.Context, r *ListRequest) (*ListReply, error) {
	reply, err := s.repository.List(ctx, &repository.ListCategoriesRequest{
		UserID: r.UserID,
		Type:   r.Type,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	return &ListReply{
		Categories: reply.Categories,
	}, nil
}

func (s *service) Update(ctx context.Context, r *UpdateRequest) (*UpdateReply, error) {
	if r.Category == nil {
		return nil, fmt.Errorf("nothing to update")
	}
	row, err := s.repository.Update(ctx, &repository.UpdateCategoryRequest{
		UserID:     r.UserID,
		CategoryID: r.CategoryID,
		Category:   r.Category,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &UpdateReply{
		Category: row,
	}, nil
}

func (s *service) Delete(ctx context.Context, r *DeleteRequest) (*DeleteReply, error) {
	_, err := s.repository.Delete(ctx, &repository.DeleteCategoriesRequest{
		UserID:      r.UserID,
		CategoryIDs: []int32{r.CategoryID},
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &DeleteReply{}, nil
}
