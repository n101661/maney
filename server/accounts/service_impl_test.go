package accounts

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
			userID      = "user-id"
			accountName = "A"
			iconID      = 11
			publicID    = "publicID"

			returnedAccountID = 9
		)
		var (
			initBalance = decimal.NewFromInt(1)
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockAccountRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().
				Create(gomock.Any(), &repository.CreateAccountsRequest{
					UserID: userID,
					Accounts: []*repository.BaseCreateAccount{
						{
							PublicID: publicID,
							BaseAccount: &repository.BaseAccount{
								Name:           accountName,
								IconID:         iconID,
								InitialBalance: initBalance,
							},
						},
					},
				}).
				Return([]*repository.Account{
					{
						ID:       returnedAccountID,
						PublicID: publicID,
						BaseAccount: &repository.BaseAccount{
							Name:           accountName,
							IconID:         iconID,
							InitialBalance: initBalance,
						},
						Balance: initBalance,
					},
				}, nil),
		)

		s, err := NewService(mockRepo, WithAccountServiceGenPublicID(func() string {
			return publicID
		}))
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Create(context.Background(), &CreateRequest{
			UserID: userID,
			Account: &BaseAccount{
				Name:           accountName,
				IconID:         iconID,
				InitialBalance: initBalance,
			},
		})
		assert.NoError(err)
		assert.Equal(&CreateReply{
			Account: &Account{
				ID:       returnedAccountID,
				PublicID: publicID,
				BaseAccount: &BaseAccount{
					Name:           accountName,
					IconID:         iconID,
					InitialBalance: initBalance,
				},
				Balance: initBalance,
			},
		}, reply)
	})
}

func Test_service_List(t *testing.T) {
	t.Run("list some of accounts", func(t *testing.T) {
		const (
			userID = "user-id"

			accountID0   = 1
			publicID0    = "publicID0"
			accountName0 = "A"
			iconID0      = 11

			accountID1   = 2
			publicID1    = "publicID1"
			accountName1 = "B"
			iconID1      = 22
		)
		var (
			initBalance0 = decimal.NewFromInt(111)
			balance0     = decimal.NewFromInt(1111)

			initBalance1 = decimal.NewFromInt(222)
			balance1     = decimal.NewFromInt(2222)
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockAccountRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().
				List(gomock.Any(), &repository.ListAccountsRequest{
					UserID: userID,
				}).
				Return(&repository.ListAccountsReply{
					Accounts: []*repository.Account{
						{
							ID:       accountID0,
							PublicID: publicID0,
							BaseAccount: &repository.BaseAccount{
								Name:           accountName0,
								IconID:         iconID0,
								InitialBalance: initBalance0,
							},
							Balance: balance0,
						},
						{
							ID:       accountID1,
							PublicID: publicID1,
							BaseAccount: &repository.BaseAccount{
								Name:           accountName1,
								IconID:         iconID1,
								InitialBalance: initBalance1,
							},
							Balance: balance1,
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
			Accounts: []*Account{
				{
					ID:       accountID0,
					PublicID: publicID0,
					BaseAccount: &BaseAccount{
						Name:           accountName0,
						IconID:         iconID0,
						InitialBalance: initBalance0,
					},
					Balance: balance0,
				},
				{
					ID:       accountID1,
					PublicID: publicID1,
					BaseAccount: &BaseAccount{
						Name:           accountName1,
						IconID:         iconID1,
						InitialBalance: initBalance1,
					},
					Balance: balance1,
				},
			},
		}, reply)
	})
	t.Run("no account", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockAccountRepository(controller)
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
			Accounts: []*Account{},
		}, reply)
	})
}

func Test_service_Update(t *testing.T) {
	t.Run("update successful", func(t *testing.T) {
		const (
			userID      = "user-id"
			accountID   = 1
			publicID    = "publicID"
			accountName = "A"
			iconID      = 11
		)
		var (
			initBalance = decimal.NewFromInt(1)
			balance     = decimal.NewFromInt(2)

			newInitBalance = decimal.NewFromInt(2)
			newBalance     = decimal.NewFromInt(3)
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockAccountRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().
				List(gomock.Any(), &repository.ListAccountsRequest{
					UserID:          userID,
					AccountPublicID: lo.ToPtr(publicID),
				}).
				Return(&repository.ListAccountsReply{
					Accounts: []*repository.Account{{
						ID:       accountID,
						PublicID: publicID,
						BaseAccount: &repository.BaseAccount{
							Name:           accountName,
							IconID:         iconID,
							InitialBalance: initBalance,
						},
						Balance: balance,
					}},
				}, nil),
			mockRepo.EXPECT().
				Update(gomock.Any(), &repository.UpdateAccountRequest{
					UserID:          userID,
					AccountPublicID: publicID,
					Account: &repository.BaseAccount{
						Name:           accountName,
						IconID:         iconID,
						InitialBalance: newInitBalance,
					},
					BalanceDelta: lo.ToPtr(newInitBalance.Sub(initBalance)),
				}).
				Return(&repository.Account{
					ID:       accountID,
					PublicID: publicID,
					BaseAccount: &repository.BaseAccount{
						Name:           accountName,
						IconID:         iconID,
						InitialBalance: newInitBalance,
					},
					Balance: newBalance,
				}, nil),
		)

		s, err := NewService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Update(context.Background(), &UpdateRequest{
			UserID:          userID,
			AccountPublicID: publicID,
			Account: &BaseAccount{
				Name:           accountName,
				IconID:         iconID,
				InitialBalance: newInitBalance,
			},
		})
		assert.NoError(err)
		assert.Equal(&UpdateReply{
			Account: &Account{
				ID:       accountID,
				PublicID: publicID,
				BaseAccount: &BaseAccount{
					Name:           accountName,
					IconID:         iconID,
					InitialBalance: newInitBalance,
				},
				Balance: newBalance,
			},
		}, reply)
	})
	t.Run("account not found", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockAccountRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, repository.ErrDataNotFound),
		)

		s, err := NewService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Update(context.Background(), &UpdateRequest{
			UserID:          "user-id",
			AccountPublicID: "1",
			Account: &BaseAccount{
				Name:           "accountName",
				IconID:         2,
				InitialBalance: decimal.Zero,
			},
		})
		assert.ErrorIs(err, ErrAccountNotFound)
		assert.Nil(reply)
	})
}

func Test_service_Delete(t *testing.T) {
	t.Run("delete successful", func(t *testing.T) {
		const (
			userID   = "user-id"
			publicID = "publicID"

			accountID   = 1
			accountName = "A"
			iconID      = 11
		)
		var (
			initBalance = decimal.NewFromInt(111)
			balance     = decimal.NewFromInt(1111)
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockAccountRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().
				Delete(gomock.Any(), &repository.DeleteAccountsRequest{
					UserID:           userID,
					AccountPublicIDs: []string{publicID},
				}).
				Return([]*repository.Account{
					{
						ID:       accountID,
						PublicID: publicID,
						BaseAccount: &repository.BaseAccount{
							Name:           accountName,
							IconID:         iconID,
							InitialBalance: initBalance,
						},
						Balance: balance,
					},
				}, nil),
		)

		s, err := NewService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Delete(context.Background(), &DeleteRequest{
			UserID:          userID,
			AccountPublicID: publicID,
		})
		assert.NoError(err)
		assert.Equal(&DeleteReply{}, reply)
	})
	t.Run("account not found", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := repository.NewMockAccountRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, repository.ErrDataNotFound),
		)

		s, err := NewService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Delete(context.Background(), &DeleteRequest{
			UserID:          "userID",
			AccountPublicID: "1",
		})
		assert.ErrorIs(err, ErrAccountNotFound)
		assert.Nil(reply)
	})
}
