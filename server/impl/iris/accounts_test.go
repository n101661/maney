package iris

import (
	"errors"
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
			user := RegisterMockLogIn(db)
			myTestServer.db = db

			resp, err := httpDoWithToken("POST", addr, strings.NewReader(s.RequestBody), MustLogIn(user))
			assert.NoError(err, s.Name)
			assert.EqualValues(iris.StatusBadRequest, resp.StatusCode, s.Name)

			db.AssertExpectations(t)
		}
	}
	{ // the account has already existed
		db := NewMockDB()
		user := RegisterMockLogIn(db)
		db.accountService.On("Create", user.ID, dbModels.AssetAccount{
			Name:           "account-name",
			IconOID:        0,
			InitialBalance: decimal.Zero,
			Balance:        decimal.Zero,
		}).Return(uint64(0), database.ErrResourceExisted).Once()
		myTestServer.db = db

		resp, err := httpDoWithToken("POST", addr, strings.NewReader(`{
			"name": "account-name",
			"icon_oid": "0",
			"initial_balance": "0"
		}`), MustLogIn(user))
		assert.NoError(err)
		assert.EqualValues(iris.StatusConflict, resp.StatusCode)

		db.AssertExpectations(t)
	}
	{ // failed to create account(unexpected error)
		db := NewMockDB()
		user := RegisterMockLogIn(db)
		db.accountService.On("Create", user.ID, dbModels.AssetAccount{
			Name:           "account-name",
			IconOID:        0,
			InitialBalance: decimal.Zero,
			Balance:        decimal.Zero,
		}).Return(uint64(0), errors.New("unexpected error")).Once()
		myTestServer.db = db

		resp, err := httpDoWithToken("POST", addr, strings.NewReader(`{
			"name": "account-name",
			"icon_oid": "0",
			"initial_balance": "0"
		}`), MustLogIn(user))
		assert.NoError(err)
		assert.EqualValues(iris.StatusInternalServerError, resp.StatusCode)

		db.AssertExpectations(t)
	}
	{ // create account successful
		db := NewMockDB()
		user := RegisterMockLogIn(db)
		db.accountService.On("Create", user.ID, dbModels.AssetAccount{
			Name:           "account-name",
			IconOID:        0,
			InitialBalance: decimal.Zero,
			Balance:        decimal.Zero,
		}).Return(uint64(0), nil).Once()
		myTestServer.db = db

		resp, err := httpDoWithToken("POST", addr, strings.NewReader(`{
			"name": "account-name",
			"icon_oid": "0",
			"initial_balance": "0"
		}`), MustLogIn(user))
		assert.NoError(err)
		assert.EqualValues(iris.StatusOK, resp.StatusCode)

		db.AssertExpectations(t)
	}
}

func TestServer_ListAccounts(t *testing.T) {
	assert := assert.New(t)

	const addr = "http://" + serverAddr + "/users/accounts"

	{ // failed to get accounts(unexpected error)
		db := NewMockDB()
		user := RegisterMockLogIn(db)
		db.accountService.On("List", user.ID).Return(([]dbModels.AssetAccount)(nil), errors.New("unexpected error")).Once()
		myTestServer.db = db

		resp, err := httpDoWithToken("GET", addr, nil, MustLogIn(user))
		assert.NoError(err)
		assert.EqualValues(iris.StatusInternalServerError, resp.StatusCode)

		db.AssertExpectations(t)
	}
	{ // get accounts successful
		db := NewMockDB()
		user := RegisterMockLogIn(db)
		db.accountService.On("List", user.ID).Return([]dbModels.AssetAccount{{
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

		resp, err := httpDoWithToken("GET", addr, nil, MustLogIn(user))
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

		db.AssertExpectations(t)
	}
}

func TestServer_UpdateAccount(t *testing.T) {
	assert := assert.New(t)

	const addr = "http://" + serverAddr + "/users/accounts"

	{ // missing account oid
		db := NewMockDB()
		user := RegisterMockLogIn(db)
		myTestServer.db = db

		resp, err := httpDoWithToken("PUT", addr+"/", strings.NewReader(`{
			"name": "my-account",
			"icon_oid": "10",
			"initial_balance": "7.89"
		}`), MustLogIn(user))
		assert.NoError(err)
		assert.EqualValues(iris.StatusBadRequest, resp.StatusCode)

		db.AssertExpectations(t)
	}
	{ // invalid account oid
		db := NewMockDB()
		user := RegisterMockLogIn(db)
		myTestServer.db = db

		resp, err := httpDoWithToken("PUT", addr+"/abc", strings.NewReader(`{
			"name": "my-account",
			"icon_oid": "10",
			"initial_balance": "7.89"
		}`), MustLogIn(user))
		assert.NoError(err)
		assert.EqualValues(iris.StatusBadRequest, resp.StatusCode)

		db.AssertExpectations(t)
	}
	{ // failed to update account(unexpected error)
		db := NewMockDB()
		user := RegisterMockLogIn(db)
		db.accountService.On("Update", user.ID, dbModels.AssetAccount{
			OID:            0,
			Name:           "my-account",
			IconOID:        10,
			InitialBalance: MustDecimalFromString("7.89"),
		}).Return(errors.New("unexpected error"))
		myTestServer.db = db

		resp, err := httpDoWithToken("PUT", addr+"/0", strings.NewReader(`{
			"name": "my-account",
			"icon_oid": "10",
			"initial_balance": "7.89"
		}`), MustLogIn(user))
		assert.NoError(err)
		assert.EqualValues(iris.StatusInternalServerError, resp.StatusCode)

		db.AssertExpectations(t)
	}
	{ // update account successful
		db := NewMockDB()
		user := RegisterMockLogIn(db)
		db.accountService.On("Update", user.ID, dbModels.AssetAccount{
			OID:            0,
			Name:           "my-account",
			IconOID:        10,
			InitialBalance: MustDecimalFromString("7.89"),
		}).Return(nil)
		myTestServer.db = db

		resp, err := httpDoWithToken("PUT", addr+"/0", strings.NewReader(`{
			"name": "my-account",
			"icon_oid": "10",
			"initial_balance": "7.89"
		}`), MustLogIn(user))
		assert.NoError(err)
		assert.EqualValues(iris.StatusOK, resp.StatusCode)

		db.AssertExpectations(t)
	}
}

func TestServer_DeleteAccount(t *testing.T) {
	assert := assert.New(t)

	const addr = "http://" + serverAddr + "/users/accounts"

	{ // missing account oid
		db := NewMockDB()
		user := RegisterMockLogIn(db)
		myTestServer.db = db

		resp, err := httpDoWithToken("DELETE", addr+"/", nil, MustLogIn(user))
		assert.NoError(err)
		assert.EqualValues(iris.StatusBadRequest, resp.StatusCode)

		db.AssertExpectations(t)
	}
	{ // invalid account oid
		db := NewMockDB()
		user := RegisterMockLogIn(db)
		myTestServer.db = db

		resp, err := httpDoWithToken("DELETE", addr+"/abc", nil, MustLogIn(user))
		assert.NoError(err)
		assert.EqualValues(iris.StatusBadRequest, resp.StatusCode)

		db.AssertExpectations(t)
	}
	{ // delete account successful
		db := NewMockDB()
		user := RegisterMockLogIn(db)
		db.accountService.On("Delete", user.ID, uint64(99)).Return(nil)
		myTestServer.db = db

		resp, err := httpDoWithToken("DELETE", addr+"/99", nil, MustLogIn(user))
		assert.NoError(err)
		assert.EqualValues(iris.StatusOK, resp.StatusCode)

		db.AssertExpectations(t)
	}
}
