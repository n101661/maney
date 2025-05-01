package fees

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/n101661/maney/server/repository"
)

func Test_service_Create(t *testing.T) {
	t.Run("create successful", func(t *testing.T) {
		const (
			userID   = "user-id"
			feeName  = "A"
			publicID = "publicID"

			returnedFeeID = 9
		)
		var (
			rate = decimal.NewFromFloat(0.1)
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockFeeRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().
				Create(gomock.Any(), &repository.CreateFeesRequest{
					UserID: userID,
					Fees: []*repository.BaseCreateFee{
						{
							PublicID: publicID,
							BaseFee: &repository.BaseFee{
								Name: feeName,
								Rate: lo.ToPtr(rate),
							},
						},
					},
				}).
				Return([]*repository.Fee{
					{
						ID:       returnedFeeID,
						PublicID: publicID,
						BaseFee: &repository.BaseFee{
							Name: feeName,
							Rate: lo.ToPtr(rate),
						},
					},
				}, nil),
		)

		s, err := NewService(mockRepo, WithFeeServiceGenPublicID(func() string {
			return publicID
		}))
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Create(context.Background(), &CreateRequest{
			UserID: userID,
			Fee: &BaseFee{
				Name: feeName,
				Rate: lo.ToPtr(rate),
			},
		})
		assert.NoError(err)
		assert.Equal(&CreateReply{
			Fee: &Fee{
				ID:       returnedFeeID,
				PublicID: publicID,
				BaseFee: &BaseFee{
					Name: feeName,
					Rate: lo.ToPtr(rate),
				},
			},
		}, reply)
	})
}

func Test_service_List(t *testing.T) {
	t.Run("list some of fees", func(t *testing.T) {
		const (
			userID = "user-id"

			feeID0    = 1
			publicID0 = "publicID0"
			feeName0  = "A"

			feeID1    = 2
			publicID1 = "publicID1"
			feeName1  = "B"
		)
		var (
			rate0 = decimal.NewFromFloat(0.1)
			rate1 = decimal.NewFromFloat(0.2)
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockFeeRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().
				List(gomock.Any(), &repository.ListFeesRequest{
					UserID: userID,
				}).
				Return(&repository.ListFeesReply{
					Fees: []*repository.Fee{
						{
							ID:       feeID0,
							PublicID: publicID0,
							BaseFee: &repository.BaseFee{
								Name: feeName0,
								Rate: lo.ToPtr(rate0),
							},
						},
						{
							ID:       feeID1,
							PublicID: publicID1,
							BaseFee: &repository.BaseFee{
								Name: feeName1,
								Rate: lo.ToPtr(rate1),
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
			Fees: []*Fee{
				{
					ID:       feeID0,
					PublicID: publicID0,
					BaseFee: &BaseFee{
						Name: feeName0,
						Rate: lo.ToPtr(rate0),
					},
				},
				{
					ID:       feeID1,
					PublicID: publicID1,
					BaseFee: &BaseFee{
						Name: feeName1,
						Rate: lo.ToPtr(rate1),
					},
				},
			},
		}, reply)
	})
	t.Run("no fee", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockFeeRepository(controller)
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
			Fees: []*Fee{},
		}, reply)
	})
}

func Test_service_Update(t *testing.T) {
	t.Run("update successful", func(t *testing.T) {
		const (
			userID   = "user-id"
			feeID    = 1
			publicID = "publicID"
			feeName  = "A"
		)
		var (
			rate = decimal.NewFromFloat(0.1)
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockFeeRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().
				Update(gomock.Any(), &repository.UpdateFeeRequest{
					UserID:      userID,
					FeePublicID: publicID,
					Fee: &repository.BaseFee{
						Name: feeName,
						Rate: lo.ToPtr(rate),
					},
				}).
				Return(&repository.Fee{
					ID:       feeID,
					PublicID: publicID,
					BaseFee: &repository.BaseFee{
						Name: feeName,
						Rate: lo.ToPtr(rate),
					},
				}, nil),
		)

		s, err := NewService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Update(context.Background(), &UpdateRequest{
			UserID:      userID,
			FeePublicID: publicID,
			Fee: &BaseFee{
				Name: feeName,
				Rate: lo.ToPtr(rate),
			},
		})
		assert.NoError(err)
		assert.Equal(&UpdateReply{
			Fee: &Fee{
				ID:       feeID,
				PublicID: publicID,
				BaseFee: &BaseFee{
					Name: feeName,
					Rate: lo.ToPtr(rate),
				},
			},
		}, reply)
	})
	t.Run("fee not found", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockFeeRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, repository.ErrDataNotFound),
		)

		s, err := NewService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Update(context.Background(), &UpdateRequest{
			UserID:      "user-id",
			FeePublicID: "1",
			Fee: &BaseFee{
				Name: "feeName",
				Rate: lo.ToPtr(decimal.NewFromFloat(0.1)),
			},
		})
		assert.ErrorIs(err, ErrFeeNotFound)
		assert.Nil(reply)
	})
}

func Test_service_Delete(t *testing.T) {
	t.Run("delete successful", func(t *testing.T) {
		const (
			userID   = "user-id"
			publicID = "publicID"

			feeID   = 1
			feeName = "A"
		)
		var (
			rate = decimal.NewFromFloat(0.1)
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockFeeRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().
				Delete(gomock.Any(), &repository.DeleteFeesRequest{
					UserID:       userID,
					FeePublicIDs: []string{publicID},
				}).
				Return([]*repository.Fee{
					{
						ID:       feeID,
						PublicID: publicID,
						BaseFee: &repository.BaseFee{
							Name: feeName,
							Rate: lo.ToPtr(rate),
						},
					},
				}, nil),
		)

		s, err := NewService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Delete(context.Background(), &DeleteRequest{
			UserID:      userID,
			FeePublicID: publicID,
		})
		assert.NoError(err)
		assert.Equal(&DeleteReply{}, reply)
	})
	t.Run("fee not found", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockFeeRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, repository.ErrDataNotFound),
		)

		s, err := NewService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Delete(context.Background(), &DeleteRequest{
			UserID:      "userID",
			FeePublicID: "1",
		})
		assert.ErrorIs(err, ErrFeeNotFound)
		assert.Nil(reply)
	})
}
