package iris

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/iris-contrib/httpexpect/v2"
	"github.com/kataras/iris/v12/httptest"
	"go.uber.org/mock/gomock"

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
	mockService := users.NewMockService(controller)

	// Set up expectations of the mock service.
	mockService.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&users.LoginReply{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil).AnyTimes()
	mockService.EXPECT().Logout(gomock.Any(), gomock.Any()).Return(&users.LogoutReply{}, nil).AnyTimes()
	mockService.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(&users.SignUpReply{}, nil).AnyTimes()
	mockService.EXPECT().RefreshAccessToken(gomock.Any(), gomock.Any()).Return(&users.RefreshAccessTokenReply{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil).AnyTimes()
	mockService.EXPECT().ValidateAccessToken(gomock.Any(), gomock.Any()).Return(&users.ValidateAccessTokenReply{}, nil).AnyTimes()
	mockService.EXPECT().UpdateConfig(gomock.Any(), gomock.Any()).Return(&users.UpdateConfigReply{}, nil).AnyTimes()
	mockService.EXPECT().GetConfig(gomock.Any(), gomock.Any()).Return(&users.GetConfigReply{}, nil).AnyTimes()

	httpExpect := httptest.New(t, NewServer(&Config{}, users.NewIrisController(mockService)).app)

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
