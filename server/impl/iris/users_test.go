package iris

import (
	"errors"
	"net/http"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/assert"

	"github.com/n101661/maney/database"
	dbModels "github.com/n101661/maney/database/models"
	"github.com/n101661/maney/server/models"
)

func TestServer_LogIn(t *testing.T) {
	assert := assert.New(t)

	const addr = "http://" + serverAddr + "/log-in"

	{ // missing required field
		scenarios := []testScenario[models.LogInRequestBody]{
			{
				Name: "missing id",
				Input: models.LogInRequestBody{
					ID:       "",
					Password: myTestPassword.Raw,
				},
			}, {
				Name: "missing password",
				Input: models.LogInRequestBody{
					ID:       "my-id",
					Password: "",
				},
			},
		}

		for _, s := range scenarios {
			resp, err := http.Post(addr, "application/json", MustHTTPBody(s.Input))
			assert.NoError(err, s.Name)
			assert.EqualValues(iris.StatusBadRequest, resp.StatusCode, s.Name)
		}
	}
	{ // no such user
		db := NewMockDB()
		db.userService.On("Get", "unknown-id").Return((*dbModels.User)(nil), nil).Once()
		myTestServer.db = db

		resp, err := http.Post(addr, "application/json", MustHTTPBody(models.LogInRequestBody{
			ID:       "unknown-id",
			Password: myTestPassword.Raw,
		}))
		assert.NoError(err)
		assert.EqualValues(iris.StatusUnauthorized, resp.StatusCode)

		db.userService.Mock.AssertExpectations(t)
	}
	{ // bad password
		db := NewMockDB()
		db.userService.On("Get", "my-id").Return(&dbModels.User{
			ID:       "my-id",
			Name:     "tester",
			Password: myTestPassword.Encrypted,
		}, nil).Once()
		myTestServer.db = db

		resp, err := http.Post(addr, "application/json", MustHTTPBody(models.LogInRequestBody{
			ID:       "my-id",
			Password: "bad-password",
		}))
		assert.NoError(err)
		assert.EqualValues(iris.StatusUnauthorized, resp.StatusCode)

		db.userService.Mock.AssertExpectations(t)
	}
	{ // log in successful
		db := NewMockDB()
		db.userService.On("Get", "my-id").Return(&dbModels.User{
			ID:       "my-id",
			Name:     "tester",
			Password: myTestPassword.Encrypted,
		}, nil).Once()
		myTestServer.db = db

		resp, err := http.Post(addr, "application/json", MustHTTPBody(models.LogInRequestBody{
			ID:       "my-id",
			Password: myTestPassword.Raw,
		}))
		assert.NoError(err)
		assert.EqualValues(iris.StatusOK, resp.StatusCode)
		assert.NotEmpty(resp.Header.Get("Set-Cookie"))

		db.userService.Mock.AssertExpectations(t)
	}
}

func TestServer_LogOut(t *testing.T) {
	assert := assert.New(t)

	resp, err := http.Post("http://"+serverAddr+"/log-out", "application/json", nil)
	assert.NoError(err)
	assert.EqualValues(iris.StatusOK, resp.StatusCode)
	assert.NotEmpty(resp.Header.Get("Set-Cookie"))
}

func TestServer_SignUp(t *testing.T) {
	assert := assert.New(t)

	const addr = "http://" + serverAddr + "/sign-up"

	{ // missing required field
		scenarios := []testScenario[models.SignUpRequestBody]{{
			Name: "missing id",
			Input: models.SignUpRequestBody{
				ID:       "",
				Name:     "tester",
				Password: myTestPassword.Raw,
			},
		}, {
			Name: "missing name",
			Input: models.SignUpRequestBody{
				ID:       "my-id",
				Name:     "",
				Password: myTestPassword.Raw,
			},
		}, {
			Name: "missing password",
			Input: models.SignUpRequestBody{
				ID:       "my-id",
				Name:     "tester",
				Password: "",
			},
		}}

		for _, s := range scenarios {
			resp, err := http.Post(addr, "application/json", MustHTTPBody(s.Input))
			assert.NoError(err, s.Name)
			assert.EqualValues(iris.StatusBadRequest, resp.StatusCode, s.Name)
		}
	}
	{ // the user has already existed
		db := NewMockDB()
		db.userService.On("Create", dbModels.User{
			ID:       "my-id",
			Name:     "tester",
			Password: []byte("my-encrypted-password"),
		}).Return(database.ErrResourceExisted).Once()
		myTestServer.db = db

		resp, err := http.Post(addr, "application/json", MustHTTPBody(models.SignUpRequestBody{
			ID:       "my-id",
			Name:     "tester",
			Password: myTestPassword.Raw,
		}))
		assert.NoError(err)
		assert.EqualValues(iris.StatusConflict, resp.StatusCode)

		db.userService.AssertExpectations(t)
	}
	{ // failed to create user(unexpected error)
		db := NewMockDB()
		db.userService.On("Create", dbModels.User{
			ID:       "my-id",
			Name:     "tester",
			Password: []byte("my-encrypted-password"),
		}).Return(errors.New("unexpected error")).Once()
		myTestServer.db = db

		resp, err := http.Post(addr, "application/json", MustHTTPBody(models.SignUpRequestBody{
			ID:       "my-id",
			Name:     "tester",
			Password: myTestPassword.Raw,
		}))
		assert.NoError(err)
		assert.EqualValues(iris.StatusInternalServerError, resp.StatusCode)

		db.userService.AssertExpectations(t)
	}
	{ // sign up successful
		db := NewMockDB()
		db.userService.On("Create", dbModels.User{
			ID:       "my-id",
			Name:     "tester",
			Password: []byte("my-encrypted-password"),
		}).Return(nil).Once()
		myTestServer.db = db

		resp, err := http.Post(addr, "application/json", MustHTTPBody(models.SignUpRequestBody{
			ID:       "my-id",
			Name:     "tester",
			Password: myTestPassword.Raw,
		}))
		assert.NoError(err)
		assert.EqualValues(iris.StatusOK, resp.StatusCode)

		db.userService.AssertExpectations(t)
	}
}

func TestServer_UpdateConfig(t *testing.T) {
	assert := assert.New(t)

	{ // failed to update config(unexpected error)
		db := NewMockDB()
		RegisterMockLogIn(db)
		db.userService.On("UpdateConfig", "my-id", dbModels.UserConfig{
			CompareItemsInDifferentShop: false,
			CompareItemsInSameShop:      true,
		}).Return(errors.New("unexpected error")).Once()
		myTestServer.db = db

		resp, err := http.DefaultClient.Do(MustHTTPRequestWithToken("PUT", "http://"+serverAddr+"/users/config", models.UserConfigRequestBody{
			CompareItemsInDifferentShop: false,
			CompareItemsInSameShop:      true,
		}, MustLogIn()))
		assert.NoError(err)
		assert.EqualValues(iris.StatusInternalServerError, resp.StatusCode)

		db.userService.AssertExpectations(t)
	}
	{ // update config successful
		db := NewMockDB()
		RegisterMockLogIn(db)
		db.userService.On("UpdateConfig", "my-id", dbModels.UserConfig{
			CompareItemsInDifferentShop: false,
			CompareItemsInSameShop:      true,
		}).Return(nil).Once()
		myTestServer.db = db

		resp, err := http.DefaultClient.Do(MustHTTPRequestWithToken("PUT", "http://"+serverAddr+"/users/config", models.UserConfigRequestBody{
			CompareItemsInDifferentShop: false,
			CompareItemsInSameShop:      true,
		}, MustLogIn()))
		assert.NoError(err)
		assert.EqualValues(iris.StatusOK, resp.StatusCode)

		db.userService.AssertExpectations(t)
	}
}

func TestServer_GetConfig(t *testing.T) {
	assert := assert.New(t)

	{ // failed to get config(unexpected error)
		db := NewMockDB()
		RegisterMockLogIn(db)
		db.userService.On("GetConfig", "my-id").Return(dbModels.UserConfig{}, errors.New("unexpected error")).Once()
		myTestServer.db = db

		resp, err := http.DefaultClient.Do(MustHTTPRequestWithToken("GET", "http://"+serverAddr+"/users/config", nil, MustLogIn()))
		assert.NoError(err)
		assert.EqualValues(iris.StatusInternalServerError, resp.StatusCode)

		db.userService.AssertExpectations(t)
	}
	{ // get config successful
		type response struct {
			CompareItemsInDifferentShop bool `json:"compare_items_in_different_shop"`
			CompareItemsInSameShop      bool `json:"compare_items_in_same_shop"`
		}

		db := NewMockDB()
		RegisterMockLogIn(db)
		db.userService.On("GetConfig", "my-id").Return(dbModels.UserConfig{
			CompareItemsInDifferentShop: true,
			CompareItemsInSameShop:      false,
		}, nil).Once()
		myTestServer.db = db

		resp, err := http.DefaultClient.Do(MustHTTPRequestWithToken("GET", "http://"+serverAddr+"/users/config", nil, MustLogIn()))
		assert.NoError(err)
		assert.EqualValues(iris.StatusOK, resp.StatusCode)
		assert.Equal(&response{
			CompareItemsInDifferentShop: true,
			CompareItemsInSameShop:      false,
		}, MustGetResponseBody[response](resp.Body))

		db.userService.AssertExpectations(t)
	}
}
