package iris

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/iris-contrib/httpexpect/v2"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/n101661/maney/database"
	dbModels "github.com/n101661/maney/database/models"
	authV2 "github.com/n101661/maney/pkg/services/auth"
	"github.com/n101661/maney/pkg/testing/aaa"
	"github.com/n101661/maney/server/models"
)

func TestServer_Login(t *testing.T) {
	type Vars struct {
		userID         string
		password       string
		accessTokenID  string
		refreshTokenID string
		httpExpect     *httpexpect.Expect
	}
	aaa := aaa.New[Vars, httpexpect.Response]()

	const path = "/login"

	t.Run("missing id of request", func(t *testing.T) {
		aaa.Arrange(func() *Vars {
			_, httpExpect := NewTest(t)
			return &Vars{httpExpect: httpExpect}
		}).Act(func(v *Vars) *httpexpect.Response {
			return v.httpExpect.POST(path).WithJSON(models.LoginRequestBody{
				ID:       "",
				Password: "password",
			}).Expect()
		}).Assert(func(v *Vars, expect *httpexpect.Response) {
			expect.Status(httptest.StatusBadRequest)
		})
	})
	t.Run("missing password of request", func(t *testing.T) {
		aaa.Arrange(func() *Vars {
			_, httpExpect := NewTest(t)
			return &Vars{httpExpect: httpExpect}
		}).Act(func(v *Vars) *httpexpect.Response {
			return v.httpExpect.POST(path).WithJSON(models.LoginRequestBody{
				ID:       "id",
				Password: "",
			}).Expect()
		}).Assert(func(v *Vars, expect *httpexpect.Response) {
			expect.Status(httptest.StatusBadRequest)
		})
	})
	t.Run("no such user", func(t *testing.T) {
		aaa.Arrange(func() *Vars {
			mock, httpExpect := NewTest(t)
			vars := &Vars{
				userID:     "wrong-id",
				password:   "password",
				httpExpect: httpExpect,
			}

			mock.auth.EXPECT().
				ValidateUser(gomock.Any(), vars.userID, vars.password).
				Return(authV2.ErrUserNotFoundOrInvalidPassword)

			return vars
		}).Act(func(v *Vars) *httpexpect.Response {
			return v.httpExpect.POST(path).WithJSON(models.LoginRequestBody{
				ID:       v.userID,
				Password: v.password,
			}).Expect()
		}).Assert(func(v *Vars, expect *httpexpect.Response) {
			expect.Status(httptest.StatusUnauthorized)
		})
	})
	t.Run("log in successful", func(t *testing.T) {
		aaa.Arrange(func() *Vars {
			mock, httpExpect := NewTest(t)
			vars := &Vars{
				userID:         "id",
				password:       "password",
				accessTokenID:  "my-access-token",
				refreshTokenID: "my-refresh-token",
				httpExpect:     httpExpect,
			}

			gomock.InOrder(
				mock.auth.EXPECT().
					ValidateUser(gomock.Any(), vars.userID, vars.password).
					Return(nil),
				mock.auth.EXPECT().
					GenerateAccessToken(gomock.Any(), &authV2.TokenClaims{UserID: vars.userID}).
					Return(vars.accessTokenID, nil),
				mock.auth.EXPECT().
					GenerateRefreshToken(gomock.Any(), &authV2.TokenClaims{UserID: vars.userID}).
					Return(&authV2.Token{
						ID: vars.refreshTokenID,
						Claims: &authV2.TokenClaims{
							UserID: vars.userID,
						},
						ExpireAfter: time.Hour,
					}, nil),
			)

			return vars
		}).Act(func(v *Vars) *httpexpect.Response {
			return v.httpExpect.POST(path).WithJSON(models.LoginRequestBody{
				ID:       v.userID,
				Password: v.password,
			}).Expect()
		}).Assert(func(v *Vars, expect *httpexpect.Response) {
			expect.Status(httptest.StatusOK)
			expect.JSON().IsEqual(models.LoginResponse{
				AccessToken: v.accessTokenID,
			})
			expect.Cookie("refreshToken").Value().NotEmpty()
			expect.Cookie("refreshToken").Path().IsEqual("/auth")
			expect.Cookie("refreshToken").HasMaxAge()
		})
	})
}

func TestServer_LogIn(t *testing.T) {
	assert := assert.New(t)

	const addr = "http://" + serverAddr + "/log-in"

	{ // missing required field
		scenarios := []testScenario[models.LoginRequestBody]{
			{
				Name: "missing id",
				RequestBody: models.LoginRequestBody{
					ID:       "",
					Password: myTestPassword.Raw,
				},
			}, {
				Name: "missing password",
				RequestBody: models.LoginRequestBody{
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

		resp, err := http.Post(addr, "application/json", MustHTTPBody(models.LoginRequestBody{
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

		resp, err := http.Post(addr, "application/json", MustHTTPBody(models.LoginRequestBody{
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

		resp, err := http.Post(addr, "application/json", MustHTTPBody(models.LoginRequestBody{
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

func (s *LogInAndDoSuite) TestServer_UpdateConfig() {
	const addr = "http://" + serverAddr + "/users/config"

	{
		s.RunTest(Test{
			Name: "failed to update config(unexpected error)",
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
		})
	}
	{
		s.RunTest(Test{
			Name: "update config successful",
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
		})
	}
}

func (s *LogInAndDoSuite) TestServer_GetConfig() {
	const addr = "http://" + serverAddr + "/users/config"

	{
		s.RunTest(Test{
			Name: "failed to get config(unexpected error)",
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
		})
	}
	{
		s.RunTest(Test{
			Name: "get config successful",
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
		})
	}
}
