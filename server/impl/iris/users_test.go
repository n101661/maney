package iris

import (
	"errors"
	"testing"
	"time"

	"github.com/iris-contrib/httpexpect/v2"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"go.uber.org/mock/gomock"

	dbModels "github.com/n101661/maney/database/models"
	"github.com/n101661/maney/pkg/models"
	authV2 "github.com/n101661/maney/pkg/services/auth"
	"github.com/n101661/maney/pkg/testing/aaa"
	httpModels "github.com/n101661/maney/server/models"
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
	var (
		nonceGenerator = func() int { return 0 }
	)

	t.Run("missing id of request", func(t *testing.T) {
		aaa.Arrange(func() *Vars {
			_, httpExpect := NewTest(t, WithNonceGenerator(nonceGenerator))
			return &Vars{httpExpect: httpExpect}
		}).Act(func(v *Vars) *httpexpect.Response {
			return v.httpExpect.POST(path).WithJSON(httpModels.LoginRequestBody{
				ID:       "",
				Password: "password",
			}).Expect()
		}).Assert(func(v *Vars, expect *httpexpect.Response) {
			expect.Status(httptest.StatusBadRequest)
		})
	})
	t.Run("missing password of request", func(t *testing.T) {
		aaa.Arrange(func() *Vars {
			_, httpExpect := NewTest(t, WithNonceGenerator(nonceGenerator))
			return &Vars{httpExpect: httpExpect}
		}).Act(func(v *Vars) *httpexpect.Response {
			return v.httpExpect.POST(path).WithJSON(httpModels.LoginRequestBody{
				ID:       "id",
				Password: "",
			}).Expect()
		}).Assert(func(v *Vars, expect *httpexpect.Response) {
			expect.Status(httptest.StatusBadRequest)
		})
	})
	t.Run("no such user", func(t *testing.T) {
		aaa.Arrange(func() *Vars {
			mock, httpExpect := NewTest(t, WithNonceGenerator(nonceGenerator))
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
			return v.httpExpect.POST(path).WithJSON(httpModels.LoginRequestBody{
				ID:       v.userID,
				Password: v.password,
			}).Expect()
		}).Assert(func(v *Vars, expect *httpexpect.Response) {
			expect.Status(httptest.StatusUnauthorized)
		})
	})
	t.Run("log in successful", func(t *testing.T) {
		aaa.Arrange(func() *Vars {
			mock, httpExpect := NewTest(t, WithNonceGenerator(nonceGenerator))
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
			return v.httpExpect.POST(path).WithJSON(httpModels.LoginRequestBody{
				ID:       v.userID,
				Password: v.password,
			}).Expect()
		}).Assert(func(v *Vars, expect *httpexpect.Response) {
			expect.Status(httptest.StatusOK)
			expect.JSON().IsEqual(httpModels.LoginResponse{
				AccessToken: v.accessTokenID,
			})

			expectRefreshTokenCookie := expect.Cookie(cookieRefreshToken)
			expectRefreshTokenCookie.Value().NotEmpty()
			expectRefreshTokenCookie.Path().IsEqual(cookiePathRefreshToken)
			expectRefreshTokenCookie.HasMaxAge()
		})
	})
}

func TestServer_Logout(t *testing.T) {
	type Vars struct {
		refreshTokenID string
		httpExpect     *httpexpect.Expect
	}
	aaa := aaa.New[Vars, httpexpect.Response]()

	const path = "/auth/logout"

	t.Run("missing refresh token", func(t *testing.T) {
		aaa.Arrange(func() *Vars {
			_, httpExpect := NewTest(t)
			return &Vars{
				httpExpect: httpExpect,
			}
		}).Act(func(v *Vars) *httpexpect.Response {
			return v.httpExpect.POST(path).Expect()
		}).Assert(func(v *Vars, expect *httpexpect.Response) {
			expect.Status(httptest.StatusOK)
			expect.Cookies().IsEmpty()
		})
	})
	t.Run("log out successful", func(t *testing.T) {
		aaa.Arrange(func() *Vars {
			mock, httpExpect := NewTest(t)
			vars := &Vars{
				refreshTokenID: "my-refresh-token",
				httpExpect:     httpExpect,
			}

			mock.auth.EXPECT().
				RevokeRefreshToken(gomock.Any(), vars.refreshTokenID).
				Return(nil)

			return vars
		}).Act(func(v *Vars) *httpexpect.Response {
			return v.httpExpect.POST(path).
				WithCookie(cookieRefreshToken, v.refreshTokenID).
				Expect()
		}).Assert(func(v *Vars, expect *httpexpect.Response) {
			expect.Status(httptest.StatusOK)

			expectRefreshTokenCookie := expect.Cookie(cookieRefreshToken)
			expectRefreshTokenCookie.Value().NotEmpty()
			expectRefreshTokenCookie.Path().IsEqual(cookiePathRefreshToken)
			expectRefreshTokenCookie.NotHasMaxAge()
		})
	})
}

func TestServer_SignUp(t *testing.T) {
	type Vars struct {
		userID     string
		password   string
		httpExpect *httpexpect.Expect
	}
	aaa := aaa.New[Vars, httpexpect.Response]()

	const path = "/sign-up"

	t.Run("missing required field of request", func(t *testing.T) {
		type Test struct {
			Name    string
			Request httpModels.SignUpRequestBody
		}

		testCases := []Test{
			{
				Name: "missing id",
				Request: httpModels.SignUpRequestBody{
					ID:       "",
					Password: "password",
				},
			},
			{
				Name: "missing password",
				Request: httpModels.SignUpRequestBody{
					ID:       "id",
					Password: "",
				},
			},
		}
		for _, c := range testCases {
			t.Run(c.Name, func(t *testing.T) {
				aaa.Arrange(func() *Vars {
					_, httpExpect := NewTest(t)
					return &Vars{
						userID:     c.Request.ID,
						password:   c.Request.Password,
						httpExpect: httpExpect,
					}
				}).Act(func(v *Vars) *httpexpect.Response {
					return v.httpExpect.POST(path).WithJSON(httpModels.SignUpRequestBody{
						ID:       v.userID,
						Password: v.password,
					}).Expect()
				}).Assert(func(v *Vars, expect *httpexpect.Response) {
					expect.Status(httptest.StatusBadRequest)
				})
			})
		}
	})
	t.Run("the user has already existed", func(t *testing.T) {
		aaa.Arrange(func() *Vars {
			mock, httpExpect := NewTest(t)
			vars := &Vars{
				userID:     "my-id",
				password:   "my-password",
				httpExpect: httpExpect,
			}

			mock.auth.EXPECT().CreateUser(gomock.Any(), &models.User{
				ID:       vars.userID,
				Password: vars.password,
			}).Return(authV2.ErrUserExists)

			return vars
		}).Act(func(v *Vars) *httpexpect.Response {
			return v.httpExpect.POST(path).WithJSON(httpModels.SignUpRequestBody{
				ID:       v.userID,
				Password: v.password,
			}).Expect()
		}).Assert(func(v *Vars, expect *httpexpect.Response) {
			expect.Status(httptest.StatusConflict)
		})
	})
	t.Run("sign up successful", func(t *testing.T) {
		aaa.Arrange(func() *Vars {
			mock, httpExpect := NewTest(t)
			vars := &Vars{
				userID:     "my-id",
				password:   "my-password",
				httpExpect: httpExpect,
			}

			mock.auth.EXPECT().CreateUser(gomock.Any(), &models.User{
				ID:       vars.userID,
				Password: vars.password,
			}).Return(nil)

			return vars
		}).Act(func(v *Vars) *httpexpect.Response {
			return v.httpExpect.POST(path).WithJSON(httpModels.SignUpRequestBody{
				ID:       v.userID,
				Password: v.password,
			}).Expect()
		}).Assert(func(v *Vars, expect *httpexpect.Response) {
			expect.Status(httptest.StatusOK)
		})
	})
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
