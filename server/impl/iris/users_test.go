package iris

import (
	"errors"
	"net/http"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

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
				RequestBody: models.LogInRequestBody{
					ID:       "",
					Password: myTestPassword.Raw,
				},
			}, {
				Name: "missing password",
				RequestBody: models.LogInRequestBody{
					ID:       "my-id",
					Password: "",
				},
			},
		}

		for _, s := range scenarios {
			resp, err := http.Post(addr, "application/json", MustHTTPBody(s.RequestBody))
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

		db.AssertExpectations(t)
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

		db.AssertExpectations(t)
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

		db.AssertExpectations(t)
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
			RequestBody: models.SignUpRequestBody{
				ID:       "",
				Name:     "tester",
				Password: myTestPassword.Raw,
			},
		}, {
			Name: "missing name",
			RequestBody: models.SignUpRequestBody{
				ID:       "my-id",
				Name:     "",
				Password: myTestPassword.Raw,
			},
		}, {
			Name: "missing password",
			RequestBody: models.SignUpRequestBody{
				ID:       "my-id",
				Name:     "tester",
				Password: "",
			},
		}}

		for _, s := range scenarios {
			resp, err := http.Post(addr, "application/json", MustHTTPBody(s.RequestBody))
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

		db.AssertExpectations(t)
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

		db.AssertExpectations(t)
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

		db.AssertExpectations(t)
	}
}

func TestServer_UpdateConfig(t *testing.T) {
	const addr = "http://" + serverAddr + "/users/config"

	{ // failed to update config(unexpected error)
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			BeforeTest: func(userID string, db *mockDB) {
				db.userService.On("UpdateConfig", userID, dbModels.UserConfig{
					CompareItemsInDifferentShop: false,
					CompareItemsInSameShop:      true,
				}).Return(errors.New("unexpected error")).Once()
			},
			HTTPRequest: HTTPRequest{
				Method: "PUT",
				URL:    addr,
				Body: `{
					"compare_items_in_different_shop": false,
					"compare_items_in_same_shop": true
				}`,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusInternalServerError,
			},
		}))
	}
	{ // update config successful
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			BeforeTest: func(userID string, db *mockDB) {
				db.userService.On("UpdateConfig", userID, dbModels.UserConfig{
					CompareItemsInDifferentShop: false,
					CompareItemsInSameShop:      true,
				}).Return(nil).Once()
			},
			HTTPRequest: HTTPRequest{
				Method: "PUT",
				URL:    addr,
				Body: `{
					"compare_items_in_different_shop": false,
					"compare_items_in_same_shop": true
				}`,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusOK,
			},
		}))
	}
}

func TestServer_GetConfig(t *testing.T) {
	const addr = "http://" + serverAddr + "/users/config"

	{ // failed to get config(unexpected error)
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			BeforeTest: func(userID string, db *mockDB) {
				db.userService.On("GetConfig", userID).Return(dbModels.UserConfig{}, errors.New("unexpected error")).Once()
			},
			HTTPRequest: HTTPRequest{
				Method: "GET",
				URL:    addr,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusInternalServerError,
			},
		}))
	}
	{ // get config successful
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			BeforeTest: func(userID string, db *mockDB) {
				db.userService.On("GetConfig", userID).Return(dbModels.UserConfig{
					CompareItemsInDifferentShop: true,
					CompareItemsInSameShop:      false,
				}, nil).Once()
			},
			HTTPRequest: HTTPRequest{
				Method: "GET",
				URL:    addr,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusOK,
				BodyJSON: `{
					"compare_items_in_different_shop": true,
					"compare_items_in_same_shop": false
				}`,
			},
		}))
	}
}
