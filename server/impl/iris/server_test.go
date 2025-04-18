package iris

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/iris-contrib/httpexpect/v2"
	"github.com/kataras/iris/v12/httptest"
	"github.com/shopspring/decimal"
	"go.uber.org/mock/gomock"

	"github.com/n101661/maney/server/accounts"
	"github.com/n101661/maney/server/models"
	"github.com/n101661/maney/server/users"
)

func TestServer(t *testing.T) {
	var (
		accessToken = &users.Token{
			ID: "access-token",
			Claims: &users.TokenClaims{
				UserID: "user-id",
				Nonce:  0,
			},
			ExpireAfter: time.Hour,
		}
		refreshToken = &users.Token{
			ID: "refresh-token",
			Claims: &users.TokenClaims{
				UserID: "user-id",
				Nonce:  0,
			},
			ExpireAfter: time.Hour,
		}
	)

	controller := gomock.NewController(t)
	userService := users.NewMockService(controller)
	accountService := accounts.NewMockService(controller)

	// Set up expectations of the mock service.
	userService.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&users.LoginReply{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil).AnyTimes()
	userService.EXPECT().Logout(gomock.Any(), gomock.Any()).Return(&users.LogoutReply{}, nil).AnyTimes()
	userService.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(&users.SignUpReply{}, nil).AnyTimes()
	userService.EXPECT().RefreshAccessToken(gomock.Any(), gomock.Any()).Return(&users.RefreshAccessTokenReply{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil).AnyTimes()
	userService.EXPECT().ValidateAccessToken(gomock.Any(), gomock.Any()).Return(&users.ValidateAccessTokenReply{}, nil).AnyTimes()
	userService.EXPECT().UpdateConfig(gomock.Any(), gomock.Any()).Return(&users.UpdateConfigReply{}, nil).AnyTimes()
	userService.EXPECT().GetConfig(gomock.Any(), gomock.Any()).Return(&users.GetConfigReply{}, nil).AnyTimes()

	accountService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&accounts.CreateReply{
		Account: &accounts.Account{
			ID: 0,
			BaseAccount: &accounts.BaseAccount{
				Name:           "A",
				IconID:         0,
				InitialBalance: decimal.Zero,
			},
			Balance: decimal.Zero,
		},
	}, nil).AnyTimes()
	accountService.EXPECT().List(gomock.Any(), gomock.Any()).Return(&accounts.ListReply{
		Accounts: []*accounts.Account{{
			ID: 0,
			BaseAccount: &accounts.BaseAccount{
				Name:           "A",
				IconID:         0,
				InitialBalance: decimal.Zero,
			},
			Balance: decimal.Zero,
		}},
	}, nil).AnyTimes()
	accountService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&accounts.UpdateReply{
		Account: &accounts.Account{
			ID: 0,
			BaseAccount: &accounts.BaseAccount{
				Name:           "A",
				IconID:         0,
				InitialBalance: decimal.Zero,
			},
			Balance: decimal.Zero,
		},
	}, nil).AnyTimes()
	accountService.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(&accounts.DeleteReply{}, nil).AnyTimes()

	httpExpect := httptest.New(t, NewServer(&Config{}, &Controllers{
		User: users.NewIrisController(userService),
	}).app)

	loginResponse := httpExpect.POST("/login").WithJSON(models.LoginRequest{
		Id:       "user-id",
		Password: "password",
	}).Expect().Status(httptest.StatusOK)

	withAuthorization, err := newWithAuthorizationHandler(loginResponse)
	if err != nil {
		t.Fatalf("failed to create withAuthorization handler: %v", err)
	}

	httpExpect.POST("/auth/logout").WithCookie(users.CookieRefreshToken, "refresh-token-id").
		Expect().Status(httptest.StatusOK)

	httpExpect.POST("/auth/refresh").WithCookie(users.CookieRefreshToken, "refresh-token-id").
		Expect().Status(httptest.StatusOK)

	httpExpect.POST("/sign-up").WithJSON(models.SignUpRequest{
		Id:       "user-id",
		Password: "password",
	}).Expect().Status(httptest.StatusOK)

	withAuthorization(httpExpect.PUT("/config")).WithJSON(models.UserConfig{
		CompareItemsInDifferentShop: true,
		CompareItemsInSameShop:      true,
	}).Expect().Status(httptest.StatusOK)

	withAuthorization(httpExpect.GET("/config")).
		Expect().Status(httptest.StatusOK)

	withAuthorization(httpExpect.POST("/accounts")).WithJSON(models.BasicAccount{
		Name:           "A",
		IconId:         0,
		InitialBalance: "0",
	}).Expect().Status(httptest.StatusOK)

	withAuthorization(httpExpect.GET("/accounts")).
		Expect().Status(httptest.StatusOK)

	withAuthorization(httpExpect.PUT("/accounts/1")).WithJSON(models.BasicAccount{
		Name:           "A",
		IconId:         0,
		InitialBalance: "0",
	}).Expect().Status(httptest.StatusOK)

	withAuthorization(httpExpect.DELETE("/accounts/1")).
		Expect().Status(httptest.StatusOK)
}

func newWithAuthorizationHandler(resp *httpexpect.Response) (func(*httpexpect.Request) *httpexpect.Request, error) {
	raw := resp.Body().Raw()
	var response models.AuthenticationResponse
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		return nil, err
	}
	authorization := fmt.Sprintf("%s %s", users.AuthType, response.AccessToken)

	return func(r *httpexpect.Request) *httpexpect.Request {
		return r.WithHeader(users.HeaderAuthorization, authorization)
	}, nil
}
