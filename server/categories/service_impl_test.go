package categories

import (
	"context"
	"testing"

	"github.com/n101661/maney/server/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_service_Create(t *testing.T) {
	t.Run("create successful", func(t *testing.T) {
		const (
			userID = "userID"
			type_  = repository.CategoryTypeExpense

			id           = 1
			categoryName = "A"
			iconID       = 11
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		repo := repository.NewMockCategoryRepository(controller)
		gomock.InOrder(
			repo.EXPECT().
				Create(gomock.Any(), &repository.CreateCategoriesRequest{
					UserID: userID,
					Type:   type_,
					Categories: []*repository.BaseCategory{
						{
							Name:   categoryName,
							IconID: iconID,
						},
					},
				}).
				Return([]*repository.Category{
					{
						ID: id,
						BaseCategory: &repository.BaseCategory{
							Name:   categoryName,
							IconID: iconID,
						},
					},
				}, nil),
		)

		s, err := NewService(repo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Create(context.Background(), &CreateRequest{
			UserID: userID,
			Type:   type_,
			Category: &BaseCategory{
				Name:   categoryName,
				IconID: iconID,
			},
		})
		assert.NoError(err)
		assert.Equal(&CreateReply{
			Type: type_,
			Category: &Category{
				ID: id,
				BaseCategory: &BaseCategory{
					Name:   categoryName,
					IconID: iconID,
				},
			},
		}, reply)
	})
}

func Test_service_List(t *testing.T) {
	t.Run("list some of categories", func(t *testing.T) {
		const (
			userID = "userID"
			type_  = repository.CategoryTypeExpense

			id0           = 1
			categoryName0 = "A"
			iconID0       = 11

			id1           = 2
			categoryName1 = "B"
			iconID1       = 22
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		repo := repository.NewMockCategoryRepository(controller)
		gomock.InOrder(
			repo.EXPECT().
				List(gomock.Any(), &repository.ListCategoriesRequest{
					UserID: userID,
					Type:   type_,
				}).
				Return(&repository.ListCategoriesReply{
					Categories: []*repository.Category{
						{
							ID: id0,
							BaseCategory: &repository.BaseCategory{
								Name:   categoryName0,
								IconID: iconID0,
							},
						},
						{
							ID: id1,
							BaseCategory: &repository.BaseCategory{
								Name:   categoryName1,
								IconID: iconID1,
							},
						},
					},
				}, nil),
		)

		s, err := NewService(repo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.List(context.Background(), &ListRequest{
			UserID: userID,
			Type:   type_,
		})
		assert.NoError(err)
		assert.Equal(&ListReply{
			Categories: []*Category{
				{
					ID: id0,
					BaseCategory: &BaseCategory{
						Name:   categoryName0,
						IconID: iconID0,
					},
				},
				{
					ID: id1,
					BaseCategory: &BaseCategory{
						Name:   categoryName1,
						IconID: iconID1,
					},
				},
			},
		}, reply)
	})
	t.Run("no category", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		repo := repository.NewMockCategoryRepository(controller)
		gomock.InOrder(
			repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, repository.ErrDataNotFound),
		)

		s, err := NewService(repo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.List(context.Background(), &ListRequest{
			UserID: "user",
			Type:   0,
		})
		assert.NoError(err)
		assert.Equal(&ListReply{
			Categories: []*Category{},
		}, reply)
	})
}

func Test_service_Update(t *testing.T) {
	t.Run("update successful", func(t *testing.T) {
		const (
			userID = "userID"
			type_  = repository.CategoryTypeExpense

			id           = 1
			categoryName = "A"
			iconID       = 11
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		repo := repository.NewMockCategoryRepository(controller)
		gomock.InOrder(
			repo.EXPECT().
				Update(gomock.Any(), &repository.UpdateCategoryRequest{
					UserID:     userID,
					CategoryID: id,
					Category: &repository.BaseCategory{
						Name:   categoryName,
						IconID: iconID,
					},
				}).
				Return(&repository.Category{
					ID: id,
					BaseCategory: &repository.BaseCategory{
						Name:   categoryName,
						IconID: iconID,
					},
				}, nil),
		)

		s, err := NewService(repo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Update(context.Background(), &UpdateRequest{
			UserID:     userID,
			CategoryID: id,
			Category: &BaseCategory{
				Name:   categoryName,
				IconID: iconID,
			},
		})
		assert.NoError(err)
		assert.Equal(&UpdateReply{
			Category: &Category{
				ID: id,
				BaseCategory: &BaseCategory{
					Name:   categoryName,
					IconID: iconID,
				},
			},
		}, reply)
	})
	t.Run("category not found", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		repo := repository.NewMockCategoryRepository(controller)
		gomock.InOrder(
			repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, repository.ErrDataNotFound),
		)

		s, err := NewService(repo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Update(context.Background(), &UpdateRequest{
			UserID:     "user",
			CategoryID: 0,
			Category: &BaseCategory{
				Name:   "name",
				IconID: 0,
			},
		})
		assert.ErrorIs(err, ErrCategoryNotFound)
		assert.Nil(reply)
	})
}

func Test_service_Delete(t *testing.T) {
	t.Run("delete successful", func(t *testing.T) {
		const (
			userID = "userID"
			type_  = repository.CategoryTypeExpense

			id           = 1
			categoryName = "A"
			iconID       = 11
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		repo := repository.NewMockCategoryRepository(controller)
		gomock.InOrder(
			repo.EXPECT().
				Delete(gomock.Any(), &repository.DeleteCategoriesRequest{
					UserID:      userID,
					CategoryIDs: []int32{id},
				}).
				Return([]*repository.Category{{
					ID: id,
					BaseCategory: &repository.BaseCategory{
						Name:   categoryName,
						IconID: iconID,
					},
				}}, nil),
		)

		s, err := NewService(repo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Delete(context.Background(), &DeleteRequest{
			UserID:     userID,
			CategoryID: id,
		})
		assert.NoError(err)
		assert.Equal(&DeleteReply{}, reply)
	})
	t.Run("category not found", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		repo := repository.NewMockCategoryRepository(controller)
		gomock.InOrder(
			repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, repository.ErrDataNotFound),
		)

		s, err := NewService(repo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Delete(context.Background(), &DeleteRequest{
			UserID:     "user",
			CategoryID: 0,
		})
		assert.ErrorIs(err, ErrCategoryNotFound)
		assert.Nil(reply)
	})
}
