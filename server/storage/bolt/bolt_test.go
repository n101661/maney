package bolt

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/n101661/maney/server/storage"
	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	const dbPath = "test.db"

	s, err := New(dbPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		assert.NoError(t, s.Close())
		assert.NoError(t, os.Remove(dbPath))
	}()

	t.Run("create the user successful", func(t *testing.T) {
		assert.NoError(t, s.CreateUser(context.Background(), &storage.User{
			ID:       "tester",
			Password: []byte("password"),
		}))
	})
	t.Run("create existing user failed", func(t *testing.T) {
		assert.ErrorIs(t, s.CreateUser(context.Background(), &storage.User{
			ID:       "tester",
			Password: []byte("password-2"),
		}), storage.ErrExists)
	})
	t.Run("get the user successful", func(t *testing.T) {
		user, err := s.GetUser(context.Background(), "tester")
		assert.NoError(t, err)
		assert.Equal(t, &storage.User{
			ID:       "tester",
			Password: []byte("password"),
		}, user)
	})
	t.Run("get non-existing user failed", func(t *testing.T) {
		user, err := s.GetUser(context.Background(), "tester-2")
		assert.Nil(t, user)
		assert.ErrorIs(t, err, storage.ErrNotFound)
	})
	t.Run("create the token successful", func(t *testing.T) {
		assert.NoError(t, s.CreateToken(context.Background(), &storage.Token{
			ID: "token",
			Claim: &storage.TokenClaims{
				UserID: "tester",
				Nonce:  0,
			},
			ExpiryTime: time.Now(),
		}))
	})
	t.Run("create existing token failed", func(t *testing.T) {
		assert.ErrorIs(t, s.CreateToken(context.Background(), &storage.Token{
			ID: "token",
			Claim: &storage.TokenClaims{
				UserID: "tester-2",
				Nonce:  0,
			},
			ExpiryTime: time.Now(),
		}), storage.ErrExists)
	})
	t.Run("get the token successful", func(t *testing.T) {
		token, err := s.GetToken(context.Background(), "token")
		assert.NoError(t, err)

		// Ignore the expiry time because it depends on the unmarshalling function.
		token.ExpiryTime = time.Time{}
		assert.Equal(t, &storage.Token{
			ID: "token",
			Claim: &storage.TokenClaims{
				UserID: "tester",
				Nonce:  0,
			},
		}, token)
	})
	t.Run("get non-existing token failed", func(t *testing.T) {
		token, err := s.GetToken(context.Background(), "token-2")
		assert.Nil(t, token)
		assert.ErrorIs(t, err, storage.ErrNotFound)
	})
	t.Run("delete the token successful", func(t *testing.T) {
		token, err := s.DeleteToken(context.Background(), "token")
		assert.NoError(t, err)

		// Ignore the expiry time because it depends on the unmarshalling function.
		token.ExpiryTime = time.Time{}
		assert.Equal(t, &storage.Token{
			ID: "token",
			Claim: &storage.TokenClaims{
				UserID: "tester",
				Nonce:  0,
			},
		}, token)
	})
	t.Run("delete non-existing token successful", func(t *testing.T) {
		token, err := s.DeleteToken(context.Background(), "token-2")
		assert.NoError(t, err)
		assert.Empty(t, token)
	})
}
