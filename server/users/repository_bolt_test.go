package users

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/n101661/maney/server/repository"
)

func Test_boltRepository(t *testing.T) {
	const dbPath = "test.db"

	s, err := NewBoltRepository(dbPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		assert.NoError(t, s.Close())
		assert.NoError(t, os.Remove(dbPath))
	}()

	const (
		userID   = "tester"
		password = "password"
	)

	t.Run("create the user successful", func(t *testing.T) {
		assert.NoError(t, s.CreateUser(context.Background(), &repository.UserModel{
			ID:       userID,
			Password: []byte(password),
			Config:   &UserConfig{},
		}))
	})
	t.Run("create existing user failed", func(t *testing.T) {
		assert.ErrorIs(t, s.CreateUser(context.Background(), &repository.UserModel{
			ID:       userID,
			Password: []byte("password-2"),
			Config:   &UserConfig{},
		}), repository.ErrDataExists)
	})
	t.Run("get the user successful", func(t *testing.T) {
		user, err := s.GetUser(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, &repository.UserModel{
			ID:       userID,
			Password: []byte(password),
			Config:   &UserConfig{},
		}, user)
	})
	t.Run("get non-existing user failed", func(t *testing.T) {
		user, err := s.GetUser(context.Background(), "not found"+userID)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, repository.ErrDataNotFound)
	})
	t.Run("update the user successful", func(t *testing.T) {
		const newPassword = "new" + password

		err := s.UpdateUser(context.Background(), &repository.UserModel{
			ID:       userID,
			Password: []byte(newPassword),
			Config: &UserConfig{
				CompareItemsInDifferentShop: true,
				CompareItemsInSameShop:      true,
			},
		})
		assert.NoError(t, err)

		updatedUser, err := s.GetUser(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, &repository.UserModel{
			ID:       userID,
			Password: []byte(newPassword),
			Config: &UserConfig{
				CompareItemsInDifferentShop: true,
				CompareItemsInSameShop:      true,
			},
		}, updatedUser)
	})
	t.Run("update non-existing user", func(t *testing.T) {
		err := s.UpdateUser(context.Background(), &repository.UserModel{
			ID:       "not found" + userID,
			Password: []byte(password),
			Config:   &UserConfig{},
		})
		assert.ErrorIs(t, err, repository.ErrDataNotFound)
	})
	t.Run("create the token successful", func(t *testing.T) {
		assert.NoError(t, s.CreateToken(context.Background(), &repository.TokenModel{
			ID: "token",
			Claim: &repository.TokenClaims{
				UserID: "tester",
				Nonce:  0,
			},
			ExpiryTime: time.Now(),
		}))
	})
	t.Run("create existing token failed", func(t *testing.T) {
		assert.ErrorIs(t, s.CreateToken(context.Background(), &repository.TokenModel{
			ID: "token",
			Claim: &repository.TokenClaims{
				UserID: "tester-2",
				Nonce:  0,
			},
			ExpiryTime: time.Now(),
		}), repository.ErrDataExists)
	})
	t.Run("get the token successful", func(t *testing.T) {
		token, err := s.GetToken(context.Background(), "token")
		assert.NoError(t, err)

		// Ignore the expiry time because it depends on the unmarshaling function.
		token.ExpiryTime = time.Time{}
		assert.Equal(t, &repository.TokenModel{
			ID: "token",
			Claim: &repository.TokenClaims{
				UserID: "tester",
				Nonce:  0,
			},
		}, token)
	})
	t.Run("get non-existing token failed", func(t *testing.T) {
		token, err := s.GetToken(context.Background(), "token-2")
		assert.Nil(t, token)
		assert.ErrorIs(t, err, repository.ErrDataNotFound)
	})
	t.Run("delete the token successful", func(t *testing.T) {
		token, err := s.DeleteToken(context.Background(), "token")
		assert.NoError(t, err)

		// Ignore the expiry time because it depends on the unmarshaling function.
		token.ExpiryTime = time.Time{}
		assert.Equal(t, &repository.TokenModel{
			ID: "token",
			Claim: &repository.TokenClaims{
				UserID: "tester",
				Nonce:  0,
			},
		}, token)
	})
	t.Run("delete non-existing token successful", func(t *testing.T) {
		token, err := s.DeleteToken(context.Background(), "token-2")
		assert.ErrorIs(t, err, repository.ErrDataNotFound)
		assert.Nil(t, token)
	})
}
