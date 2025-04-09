package iris

import (
	"testing"
	"time"

	"github.com/kataras/iris/v12/httptest"
	"go.uber.org/mock/gomock"

	"github.com/n101661/maney/server/models"
	"github.com/n101661/maney/server/users"
)

func TestServer(t *testing.T) {
	controller := gomock.NewController(t)
	mockService := users.NewMockService(controller)

	// Set up expectations of the mock service.
	mockService.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&users.LoginReply{
		AccessToken: &users.Token{
			ID: "access-token",
			Claims: &users.TokenClaims{
				UserID: "user-id",
				Nonce:  0,
			},
			ExpireAfter: time.Hour,
		},
		RefreshToken: &users.Token{
			ID: "refresh-token",
			Claims: &users.TokenClaims{
				UserID: "user-id",
				Nonce:  0,
			},
			ExpireAfter: time.Hour,
		},
	}, nil).AnyTimes()
	mockService.EXPECT().Logout(gomock.Any(), gomock.Any()).Return(&users.LogoutReply{}, nil).AnyTimes()
	mockService.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(&users.SignUpReply{}, nil).AnyTimes()

	httpExpect := httptest.New(t, NewServer(&Config{}, users.NewIrisController(mockService)).app)

	httpExpect.POST("/login").WithJSON(models.LoginRequest{
		Id:       "user-id",
		Password: "password",
	}).Expect().Status(httptest.StatusOK)

	httpExpect.POST("/auth/logout").WithCookie(users.CookieRefreshToken, "refresh-token-id").
		Expect().Status(httptest.StatusOK)

	httpExpect.POST("/sign-up").WithJSON(models.SignUpRequest{
		Id:       "user-id",
		Password: "password",
	}).Expect().Status(httptest.StatusOK)
}
