package users

import (
	"errors"
	"net/http"
	"strings"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"

	"github.com/n101661/maney/pkg/utils"
	httpModels "github.com/n101661/maney/server/models"
)

const (
	HeaderAuthorization = "Authorization"
)

const (
	AuthType = "Bearer"
)

const (
	CookieRefreshToken = "refreshToken"
)

const (
	cookiePathRefreshToken = "/auth"
)

type IrisController struct {
	s Service

	opts *irisControllerOptions
}

func NewIrisController(s Service, opts ...utils.Option[irisControllerOptions]) *IrisController {
	return &IrisController{
		s:    s,
		opts: utils.ApplyOptions(&irisControllerOptions{}, opts),
	}
}

func (controller *IrisController) Login(c iris.Context) {
	var r httpModels.LoginRequest
	if err := c.ReadJSON(&r); err != nil {
		c.StatusCode(iris.StatusBadRequest)
		c.WriteString(err.Error())
		return
	}

	if r.Id == "" || r.Password == "" {
		c.StatusCode(iris.StatusBadRequest)
		return
	}

	ctx := c.Request().Context()

	reply, err := controller.s.Login(ctx, &LoginRequest{
		UserID:   r.Id,
		Password: r.Password,
	})
	if err != nil {
		if errors.Is(err, ErrUserNotFoundOrInvalidPassword) {
			c.StatusCode(iris.StatusUnauthorized)
			return
		}
		c.StatusCode(iris.StatusInternalServerError)
		c.WriteString(err.Error())
		return
	}

	c.SetCookieKV(
		CookieRefreshToken, reply.RefreshToken.ID,
		iris.CookiePath(cookiePathRefreshToken),
		iris.CookieExpires(reply.RefreshToken.ExpireAfter),
		iris.CookieHTTPOnly(true),
		iris.CookieSameSite(http.SameSiteStrictMode),
	)

	c.StatusCode(iris.StatusOK)
	err = c.JSON(&httpModels.AuthenticationResponse{
		AccessToken: reply.AccessToken.ID,
	})
	if err != nil && controller.opts.logger != nil {
		controller.opts.logger.Warnf("failed to response of Login: %v", err)
	}
}

func (controller *IrisController) Logout(c iris.Context) {
	cookie, err := c.Request().Cookie(CookieRefreshToken)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			c.StatusCode(iris.StatusOK)
			return
		}
		c.StatusCode(iris.StatusInternalServerError)
		return
	}

	_, err = controller.s.Logout(c.Request().Context(), &LogoutRequest{
		RefreshTokenID: cookie.Value,
	})
	if err != nil {
		if !(errors.Is(err, ErrInvalidToken) || errors.Is(err, ErrTokenExpired)) {
			c.StatusCode(iris.StatusInternalServerError)
			return
		}

		if controller.opts.logger != nil {
			controller.opts.logger.Warnf("receive unexpected token[%s] when revoking: %v", cookie.Value, err)
		}
	}

	c.SetCookieKV(
		CookieRefreshToken, cookie.Value,
		iris.CookiePath(cookiePathRefreshToken),
		iris.CookieExpires(0),
		iris.CookieHTTPOnly(true),
		iris.CookieSameSite(http.SameSiteStrictMode),
	)

	c.StatusCode(iris.StatusOK)
}

func (controller *IrisController) SignUp(c iris.Context) {
	var r httpModels.SignUpRequest
	if err := c.ReadJSON(&r); err != nil {
		c.StatusCode(iris.StatusBadRequest)
		c.WriteString(err.Error())
		return
	}

	if r.Id == "" || r.Password == "" {
		c.StatusCode(iris.StatusBadRequest)
		return
	}

	ctx := c.Request().Context()

	_, err := controller.s.SignUp(ctx, &SignUpRequest{
		UserID:   r.Id,
		Password: r.Password,
	})
	if err != nil {
		if errors.Is(err, ErrUserExists) {
			c.StatusCode(iris.StatusConflict)
			c.WriteString("the user id has existed")
			return
		}
		c.StatusCode(iris.StatusInternalServerError)
		return
	}

	c.StatusCode(iris.StatusOK)
}

func (controller *IrisController) ValidateAccessToken(c iris.Context) {
	accessToken := controller.getAccessToken(c)
	if accessToken == "" {
		c.StopWithStatus(iris.StatusUnauthorized)
		return
	}

	tokenReply, err := controller.s.ValidateAccessToken(c.Request().Context(), &ValidateAccessTokenRequest{
		TokenID: accessToken,
	})
	if err != nil {
		if errors.Is(err, ErrInvalidToken) || errors.Is(err, ErrTokenExpired) {
			c.StopWithStatus(iris.StatusUnauthorized)
		} else {
			c.StopWithError(iris.StatusInternalServerError, err)
		}
		return
	}

	err = c.SetUser(&user{
		Token: accessToken,
		ID:    tokenReply.UserID,
	})
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		c.WriteString(err.Error())
	}
}

func (controller *IrisController) getAccessToken(c iris.Context) string {
	h := c.GetHeader(HeaderAuthorization)
	if h == "" {
		return ""
	}

	frags := strings.SplitN(h, " ", 2)
	if frags[0] != AuthType || len(frags) != 2 {
		return ""
	}

	return frags[1]
}

func (controller *IrisController) UpdateUserConfig(c iris.Context) {
	var r httpModels.UserConfig
	if err := c.ReadJSON(&r); err != nil {
		c.StatusCode(iris.StatusBadRequest)
		c.WriteString(err.Error())
		return
	}

	userID, err := c.User().GetID()
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		c.WriteString(err.Error())
		return
	}

	_, err = controller.s.UpdateConfig(c.Request().Context(), &UpdateConfigRequest{
		UserID: userID,
		Config: &r,
	})
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			c.StatusCode(iris.StatusBadRequest)
			return
		}
		c.StatusCode(iris.StatusInternalServerError)
		c.WriteString(err.Error())
		return
	}

	c.StatusCode(iris.StatusOK)
}

func (controller *IrisController) GetUserConfig(c iris.Context) {
	userID, err := c.User().GetID()
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		c.WriteString(err.Error())
		return
	}

	reply, err := controller.s.GetConfig(c.Request().Context(), &GetConfigRequest{
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			c.StatusCode(iris.StatusBadRequest)
			return
		}
		c.StatusCode(iris.StatusInternalServerError)
		c.WriteString(err.Error())
		return
	}

	c.StatusCode(iris.StatusOK)
	c.JSON(reply.Data)
}

type irisControllerOptions struct {
	logger *golog.Logger
}

func WithLogger(logger *golog.Logger) utils.Option[irisControllerOptions] {
	return func(o *irisControllerOptions) {
		o.logger = logger
	}
}

type user struct {
	Token string
	ID    string
}

func (u *user) GetRaw() (interface{}, error) {
	return u, nil
}

func (u *user) GetAuthorization() (string, error) {
	return AuthType, nil
}

func (u *user) GetID() (string, error) {
	return u.ID, nil
}

func (u *user) GetToken() ([]byte, error) {
	return []byte(u.Token), nil
}
