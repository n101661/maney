package shops

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/n101661/maney/server/repository"
)

func Test_service_Create(t *testing.T) {
	t.Run("create successful", func(t *testing.T) {
		const (
			userID   = "user-id"
			shopName = "A"
			address  = "address"
			publicID = "publicID"

			returnedShopID = 9
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockShopRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().
				Create(gomock.Any(), &repository.CreateShopsRequest{
					UserID: userID,
					Shops: []*repository.BaseCreateShop{
						{
							PublicID: publicID,
							BaseShop: &repository.BaseShop{
								Name:    shopName,
								Address: address,
							},
						},
					},
				}).
				Return([]*repository.Shop{
					{
						ID:       returnedShopID,
						PublicID: publicID,
						BaseShop: &repository.BaseShop{
							Name:    shopName,
							Address: address,
						},
					},
				}, nil),
		)

		s, err := NewService(mockRepo, WithShopServiceGenPublicID(func() string {
			return publicID
		}))
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Create(context.Background(), &CreateRequest{
			UserID: userID,
			Shop: &BaseShop{
				Name:    shopName,
				Address: address,
			},
		})
		assert.NoError(err)
		assert.Equal(&CreateReply{
			Shop: &Shop{
				ID:       returnedShopID,
				PublicID: publicID,
				BaseShop: &BaseShop{
					Name:    shopName,
					Address: address,
				},
			},
		}, reply)
	})
}

func Test_service_List(t *testing.T) {
	t.Run("list some of shops", func(t *testing.T) {
		const (
			userID = "user-id"

			shopID0   = 1
			publicID0 = "publicID0"
			shopName0 = "A"
			address0  = "address0"

			shopID1   = 2
			publicID1 = "publicID1"
			shopName1 = "B"
			address1  = "address1"
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockShopRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().
				List(gomock.Any(), &repository.ListShopsRequest{
					UserID: userID,
				}).
				Return(&repository.ListShopsReply{
					Shops: []*repository.Shop{
						{
							ID:       shopID0,
							PublicID: publicID0,
							BaseShop: &repository.BaseShop{
								Name:    shopName0,
								Address: address0,
							},
						},
						{
							ID:       shopID1,
							PublicID: publicID1,
							BaseShop: &repository.BaseShop{
								Name:    shopName1,
								Address: address1,
							},
						},
					},
				}, nil),
		)

		s, err := NewService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.List(context.Background(), &ListRequest{
			UserID: userID,
		})
		assert.NoError(err)
		assert.Equal(&ListReply{
			Shops: []*Shop{
				{
					ID:       shopID0,
					PublicID: publicID0,
					BaseShop: &BaseShop{
						Name:    shopName0,
						Address: address0,
					},
				},
				{
					ID:       shopID1,
					PublicID: publicID1,
					BaseShop: &BaseShop{
						Name:    shopName1,
						Address: address1,
					},
				},
			},
		}, reply)
	})
	t.Run("no shop", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockShopRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, repository.ErrDataNotFound),
		)

		s, err := NewService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.List(context.Background(), &ListRequest{
			UserID: "user-id",
		})
		assert.NoError(err)
		assert.Equal(&ListReply{
			Shops: []*Shop{},
		}, reply)
	})
}

func Test_service_Update(t *testing.T) {
	t.Run("update successful", func(t *testing.T) {
		const (
			userID   = "user-id"
			shopID   = 1
			publicID = "publicID"
			shopName = "A"
			address  = "address"
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockShopRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().
				Update(gomock.Any(), &repository.UpdateShopRequest{
					UserID:       userID,
					ShopPublicID: publicID,
					Shop: &repository.BaseShop{
						Name:    shopName,
						Address: address,
					},
				}).
				Return(&repository.Shop{
					ID:       shopID,
					PublicID: publicID,
					BaseShop: &repository.BaseShop{
						Name:    shopName,
						Address: address,
					},
				}, nil),
		)

		s, err := NewService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Update(context.Background(), &UpdateRequest{
			UserID:       userID,
			ShopPublicID: publicID,
			Shop: &BaseShop{
				Name:    shopName,
				Address: address,
			},
		})
		assert.NoError(err)
		assert.Equal(&UpdateReply{
			Shop: &Shop{
				ID:       shopID,
				PublicID: publicID,
				BaseShop: &BaseShop{
					Name:    shopName,
					Address: address,
				},
			},
		}, reply)
	})
	t.Run("shop not found", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockShopRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, repository.ErrDataNotFound),
		)

		s, err := NewService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Update(context.Background(), &UpdateRequest{
			UserID:       "user-id",
			ShopPublicID: "1",
			Shop: &BaseShop{
				Name:    "shopName",
				Address: "Address",
			},
		})
		assert.ErrorIs(err, ErrShopNotFound)
		assert.Nil(reply)
	})
}

func Test_service_Delete(t *testing.T) {
	t.Run("delete successful", func(t *testing.T) {
		const (
			userID   = "user-id"
			publicID = "publicID"

			shopID   = 1
			shopName = "A"
			address  = "address"
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockShopRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().
				Delete(gomock.Any(), &repository.DeleteShopsRequest{
					UserID:        userID,
					ShopPublicIDs: []string{publicID},
				}).
				Return([]*repository.Shop{
					{
						ID:       shopID,
						PublicID: publicID,
						BaseShop: &repository.BaseShop{
							Name:    shopName,
							Address: address,
						},
					},
				}, nil),
		)

		s, err := NewService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Delete(context.Background(), &DeleteRequest{
			UserID:       userID,
			ShopPublicID: publicID,
		})
		assert.NoError(err)
		assert.Equal(&DeleteReply{}, reply)
	})
	t.Run("shop not found", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockShopRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, repository.ErrDataNotFound),
		)

		s, err := NewService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Delete(context.Background(), &DeleteRequest{
			UserID:       "userID",
			ShopPublicID: "1",
		})
		assert.ErrorIs(err, ErrShopNotFound)
		assert.Nil(reply)
	})
}
