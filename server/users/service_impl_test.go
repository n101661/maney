package users

import (
	"context"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_service_Login(t *testing.T) {
	t.Run("login successful", func(t *testing.T) {
		const (
			userID   = "id"
			password = "password"
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := NewMockRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(&UserModel{
				ID:       userID,
				Password: lo.Must(encryptPassword(password, defaultOptions().saltPasswordRound)),
			}, nil),
			mockRepo.EXPECT().CreateToken(gomock.Any(), gomock.Any()).Return(nil),
		)

		s, err := newService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Login(context.Background(), &LoginRequest{
			UserID:   userID,
			Password: password,
		})
		assert.NoError(err)
		assert.NotNil(reply)
	})
	t.Run("no such user", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := NewMockRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(nil, ErrDataNotFound),
		)

		s, err := newService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Login(context.Background(), &LoginRequest{
			UserID:   "id",
			Password: "password",
		})
		assert.ErrorIs(err, ErrUserNotFoundOrInvalidPassword)
		assert.Nil(reply)
	})
	t.Run("invalid password", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := NewMockRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(&UserModel{
				ID:       "id",
				Password: lo.Must(encryptPassword("password", defaultOptions().saltPasswordRound)),
			}, nil),
		)

		s, err := newService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Login(context.Background(), &LoginRequest{
			UserID:   "id",
			Password: "wrong-password",
		})
		assert.ErrorIs(err, ErrUserNotFoundOrInvalidPassword)
		assert.Nil(reply)
	})
}

func Test_service_Logout(t *testing.T) {
	t.Run("logout successful", func(t *testing.T) {
		const (
			tokenID = "refresh-token-id"
		)

		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := NewMockRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().DeleteToken(gomock.Any(), gomock.Any()).Return(&TokenModel{
				ID: tokenID,
				Claim: &TokenClaims{
					UserID: "user-id",
					Nonce:  0,
				},
				ExpiryTime: time.Now().Add(time.Hour),
			}, nil),
		)

		s, err := newService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Logout(context.Background(), &LogoutRequest{
			RefreshTokenID: tokenID,
		})
		assert.NoError(err)
		assert.Equal(&LogoutReply{}, reply)
	})
	t.Run("token not found", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := NewMockRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().DeleteToken(gomock.Any(), gomock.Any()).Return(nil, ErrDataNotFound),
		)

		s, err := newService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Logout(context.Background(), &LogoutRequest{
			RefreshTokenID: "refresh-token-id",
		})
		assert.ErrorIs(err, ErrInvalidToken)
		assert.Equal(&LogoutReply{}, reply)
	})
	t.Run("token is expired", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := NewMockRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().DeleteToken(gomock.Any(), gomock.Any()).Return(&TokenModel{
				ID: "token",
				Claim: &TokenClaims{
					UserID: "user-id",
					Nonce:  0,
				},
				ExpiryTime: time.Now().Add(-time.Hour),
			}, nil),
		)

		s, err := newService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Logout(context.Background(), &LogoutRequest{
			RefreshTokenID: "refresh-token-id",
		})
		assert.ErrorIs(err, ErrTokenExpired)
		assert.Equal(&LogoutReply{}, reply)
	})
}

func Test_service_SignUp(t *testing.T) {
	t.Run("sign up successful", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := NewMockRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil),
		)

		s, err := newService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.SignUp(context.Background(), &SignUpRequest{
			UserID:   "id",
			Password: "password",
		})
		assert.NoError(err)
		assert.Equal(&SignUpReply{}, reply)
	})
	t.Run("user already exists", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := NewMockRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(ErrDataExists),
		)

		s, err := newService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.SignUp(context.Background(), &SignUpRequest{
			UserID:   "id",
			Password: "password",
		})
		assert.ErrorIs(err, ErrUserExists)
		assert.Nil(reply)
	})
}

func Test_service_UpdateConfig(t *testing.T) {
	t.Run("update config successful", func(t *testing.T) {
		assert := assert.New(t)

		const (
			userID = "user-id"
		)

		controller := gomock.NewController(t)
		mockRepo := NewMockRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().UpdateUser(gomock.Any(), &UserModel{
				ID: userID,
				Config: &UserConfig{
					CompareItemsInDifferentShop: false,
					CompareItemsInSameShop:      true,
				},
			}).Return(nil),
		)

		s, err := newService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.UpdateConfig(context.Background(), &UpdateConfigRequest{
			UserID: userID,
			Config: &UserConfig{
				CompareItemsInDifferentShop: false,
				CompareItemsInSameShop:      true,
			},
		})
		assert.NoError(err)
		assert.Equal(&UpdateConfigReply{}, reply)
	})
	t.Run("user not found", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := NewMockRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(ErrDataNotFound),
		)

		s, err := newService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.UpdateConfig(context.Background(), &UpdateConfigRequest{
			UserID: "user-id",
			Config: &UserConfig{
				CompareItemsInDifferentShop: false,
				CompareItemsInSameShop:      true,
			},
		})
		assert.ErrorIs(err, ErrResourceNotFound)
		assert.Nil(reply)
	})
}

func Test_service_GetConfig(t *testing.T) {
	t.Run("get config successful", func(t *testing.T) {
		assert := assert.New(t)

		const (
			userID = "user-id"
		)

		controller := gomock.NewController(t)
		mockRepo := NewMockRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().GetUser(gomock.Any(), userID).Return(&UserModel{
				ID:       userID,
				Password: []byte("password"),
				Config: &UserConfig{
					CompareItemsInDifferentShop: false,
					CompareItemsInSameShop:      true,
				},
			}, nil),
		)

		s, err := newService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.GetConfig(context.Background(), &GetConfigRequest{
			UserID: userID,
		})
		assert.NoError(err)
		assert.Equal(&GetConfigReply{
			Data: &UserConfig{
				CompareItemsInDifferentShop: false,
				CompareItemsInSameShop:      true,
			},
		}, reply)
	})
	t.Run("user not found", func(t *testing.T) {
		assert := assert.New(t)

		controller := gomock.NewController(t)
		mockRepo := NewMockRepository(controller)
		gomock.InOrder(
			mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(nil, ErrDataNotFound),
		)

		s, err := newService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.GetConfig(context.Background(), &GetConfigRequest{
			UserID: "user-id",
		})
		assert.ErrorIs(err, ErrResourceNotFound)
		assert.Nil(reply)
	})
}

func newService(repo Repository) (Service, error) {
	return NewService(
		repo,
		[]byte("access-token-signing-key"),
		[]byte("refresh-token-signing-key"),
	)
}
