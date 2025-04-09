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
	headerAuthorization = "Authorization"
)

const (
	authType = "Bearer"
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
		c.StatusCode(iris.StatusInternalServerError)
		return
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
	h := c.GetHeader(headerAuthorization)
	if h == "" {
		c.StopWithStatus(iris.StatusUnauthorized)
		return
	}

	frags := strings.SplitN(h, " ", 2)
	if frags[0] != authType || len(frags) != 2 {
		c.StopWithStatus(iris.StatusUnauthorized)
		return
	}

	_, err := controller.s.ValidateAccessToken(c.Request().Context(), &ValidateAccessTokenRequest{
		TokenID: frags[1],
	})
	if err != nil {
		if errors.Is(err, ErrInvalidToken) || errors.Is(err, ErrTokenExpired) {
			c.StopWithStatus(iris.StatusUnauthorized)
		} else {
			c.StopWithError(iris.StatusInternalServerError, err)
		}
		return
	}
}

type irisControllerOptions struct {
	logger *golog.Logger
}

func WithLogger(logger *golog.Logger) utils.Option[irisControllerOptions] {
	return func(o *irisControllerOptions) {
		o.logger = logger
	}
}
