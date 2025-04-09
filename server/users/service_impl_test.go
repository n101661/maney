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
			mockRepo.EXPECT().DeleteToken(gomock.Any(), gomock.Any()).Return(&TokenModel{}, nil),
		)

		s, err := newService(mockRepo)
		if err != nil {
			t.Fatal(err)
		}

		reply, err := s.Logout(context.Background(), &LogoutRequest{
			RefreshTokenID: "refresh-token-id",
		})
		assert.NoError(err)
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

func newService(repo Repository) (Service, error) {
	return NewService(
		repo,
		[]byte("access-token-signing-key"),
		[]byte("refresh-token-signing-key"),
	)
}
