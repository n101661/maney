package iris

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/n101661/maney/database"
	dbModels "github.com/n101661/maney/database/models"
)

func TestServer_CreateAccount(t *testing.T) {
	assert := assert.New(t)

	const addr = "http://" + serverAddr + "/users/accounts"

	{ // missing required field
		scenarios := []testScenario[string]{{
			Name: "missing name",
			RequestBody: `{
				"name": "",
				"icon_oid": "0",
				"initial_balance": "0"
			}`,
		}}

		for _, s := range scenarios {
			db := NewMockDB()
			RegisterMockLogIn(db)
			myTestServer.db = db

			resp, err := http.DefaultClient.Do(MustHTTPRequestWithToken("POST", addr, strings.NewReader(s.RequestBody), MustLogIn()))
			assert.NoError(err, s.Name)
			assert.EqualValues(iris.StatusBadRequest, resp.StatusCode, s.Name)
		}
	}
	{ // the account has already existed
		db := NewMockDB()
		RegisterMockLogIn(db)
		db.accountService.On("Create", "my-id", dbModels.AssetAccount{
			Name:           "account-name",
			IconOID:        0,
			InitialBalance: decimal.Zero,
			Balance:        decimal.Zero,
		}).Return(uint64(0), database.ErrResourceExisted).Once()
		myTestServer.db = db

		resp, err := http.DefaultClient.Do(MustHTTPRequestWithToken("POST", addr, strings.NewReader(`{
			"name": "account-name",
			"icon_oid": "0",
			"initial_balance": "0"
		}`), MustLogIn()))
		assert.NoError(err)
		assert.EqualValues(iris.StatusConflict, resp.StatusCode)

		db.accountService.AssertExpectations(t)
	}
	{ // failed to create account(unexpected error)
		db := NewMockDB()
		RegisterMockLogIn(db)
		db.accountService.On("Create", "my-id", dbModels.AssetAccount{
			Name:           "account-name",
			IconOID:        0,
			InitialBalance: decimal.Zero,
			Balance:        decimal.Zero,
		}).Return(uint64(0), errors.New("unexpected error")).Once()
		myTestServer.db = db

		resp, err := http.DefaultClient.Do(MustHTTPRequestWithToken("POST", addr, strings.NewReader(`{
			"name": "account-name",
			"icon_oid": "0",
			"initial_balance": "0"
		}`), MustLogIn()))
		assert.NoError(err)
		assert.EqualValues(iris.StatusInternalServerError, resp.StatusCode)

		db.accountService.AssertExpectations(t)
	}
	{ // create account successful
		db := NewMockDB()
		RegisterMockLogIn(db)
		db.accountService.On("Create", "my-id", dbModels.AssetAccount{
			Name:           "account-name",
			IconOID:        0,
			InitialBalance: decimal.Zero,
			Balance:        decimal.Zero,
		}).Return(uint64(0), nil).Once()
		myTestServer.db = db

		resp, err := http.DefaultClient.Do(MustHTTPRequestWithToken("POST", addr, strings.NewReader(`{
			"name": "account-name",
			"icon_oid": "0",
			"initial_balance": "0"
		}`), MustLogIn()))
		assert.NoError(err)
		assert.EqualValues(iris.StatusOK, resp.StatusCode)

		db.accountService.AssertExpectations(t)
	}
}

func TestServer_ListAccounts(t *testing.T) {
	assert := assert.New(t)

	const addr = "http://" + serverAddr + "/users/accounts"

	{ // failed to get accounts(unexpected error)
		db := NewMockDB()
		RegisterMockLogIn(db)
		db.accountService.On("List", "my-id").Return(([]dbModels.AssetAccount)(nil), errors.New("unexpected error")).Once()
		myTestServer.db = db

		resp, err := http.DefaultClient.Do(MustHTTPRequestWithToken("GET", addr, nil, MustLogIn()))
		assert.NoError(err)
		assert.EqualValues(iris.StatusInternalServerError, resp.StatusCode)

		db.accountService.AssertExpectations(t)
	}
	{ // get accounts successful
		db := NewMockDB()
		RegisterMockLogIn(db)
		db.accountService.On("List", "my-id").Return([]dbModels.AssetAccount{{
			OID:            0,
			Name:           "account-1",
			IconOID:        0,
			InitialBalance: decimal.Zero,
			Balance:        decimal.NewFromFloat32(12.3),
		}, {
			OID:            1,
			Name:           "account-2",
			IconOID:        2,
			InitialBalance: decimal.NewFromFloat32(12.3),
			Balance:        decimal.NewFromFloat32(4.56),
		}}, nil).Once()
		myTestServer.db = db

		resp, err := http.DefaultClient.Do(MustHTTPRequestWithToken("GET", addr, nil, MustLogIn()))
		assert.NoError(err)
		assert.EqualValues(iris.StatusOK, resp.StatusCode)
		assert.JSONEq(`[
			{
				"oid": "0",
				"name": "account-1",
				"icon_oid": "0",
				"initial_balance": "0",
				"balance": "12.3"
			}, {
				"oid": "1",
				"name": "account-2",
				"icon_oid": "2",
				"initial_balance": "12.3",
				"balance": "4.56"
			}
		]`, MustGetResponseBodyJSON(resp.Body))

		db.accountService.AssertExpectations(t)
	}
}

func TestServer_UpdateAccount(t *testing.T) {
	// assert := assert.New(t)

	{ // failed to update account(unexpected error)

	}
	{ // update account successful

	}
}

func TestServer_DeleteAccount(t *testing.T) {
	// assert := assert.New(t)

}
