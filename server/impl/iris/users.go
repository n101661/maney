package iris

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/kataras/iris/v12"

	"github.com/n101661/maney/database"
	dbModels "github.com/n101661/maney/database/models"
	authV2 "github.com/n101661/maney/pkg/services/auth"
	httpModels "github.com/n101661/maney/server/models"
)

func (s *Server) Login(c iris.Context) {
	var r httpModels.LoginRequestBody
	if err := c.ReadJSON(&r); err != nil {
		c.StatusCode(iris.StatusBadRequest)
		c.WriteString(err.Error())
		return
	}

	ctx := c.Request().Context()

	err := s.authService.ValidateUser(ctx, r.ID, r.Password)
	if err != nil {
		if errors.Is(err, authV2.ErrUserNotFoundOrInvalidPassword) {
			c.StatusCode(iris.StatusUnauthorized)
			return
		}
		c.StatusCode(iris.StatusInternalServerError)
		c.WriteString(err.Error())
		return
	}

	accessToken, refreshToken, err := s.generateToken(ctx, &authV2.TokenClaims{
		UserID: r.ID,
	})
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		c.WriteString(err.Error())
		return
	}

	c.SetCookieKV(
		cookieRefreshToken, refreshToken.ID,
		iris.CookiePath(cookiePathRefreshToken),
		iris.CookieExpires(refreshToken.ExpireAfter),
		iris.CookieHTTPOnly(true),
		iris.CookieSameSite(http.SameSiteStrictMode),
	)

	c.StatusCode(iris.StatusOK)
	err = c.JSON(&httpModels.LoginResponse{
		AccessToken: accessToken,
	})
	if err != nil {
		log.Printf("failed to response: %v\n", err)
	}
}

func (s *Server) generateToken(
	ctx context.Context,
	claims *authV2.TokenClaims,
) (accessTokenID string, refreshToken *authV2.Token, err error) {
	accessTokenID, err = s.authService.GenerateAccessToken(ctx, claims)
	if err != nil {
		return "", nil, err
	}

	refreshToken, err = s.authService.GenerateRefreshToken(ctx, claims)
	if err != nil {
		return "", nil, err
	}

	return accessTokenID, refreshToken, nil
}

func (s *Server) Logout(c iris.Context) {
	cookie, err := c.Request().Cookie(cookieRefreshToken)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			c.StatusCode(iris.StatusOK)
			return
		}
		c.StatusCode(iris.StatusInternalServerError)
		return
	}

	err = s.authService.RevokeRefreshToken(c.Request().Context(), cookie.Value)
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		return
	}

	c.SetCookieKV(
		cookieRefreshToken, cookie.Value,
		iris.CookiePath(cookiePathRefreshToken),
		iris.CookieExpires(0),
		iris.CookieHTTPOnly(true),
		iris.CookieSameSite(http.SameSiteStrictMode),
	)

	c.StatusCode(iris.StatusOK)
}

func (s *Server) SignUp(c iris.Context) {
	var r httpModels.SignUpRequestBody
	if err := c.ReadJSON(&r); err != nil {
		c.StatusCode(iris.StatusBadRequest)
		c.WriteString(err.Error())
		return
	}

	pwd, err := s.auth.EncryptPassword(r.Password)
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		c.WriteString(err.Error())
		return
	}

	err = s.db.User().Create(dbModels.User{
		ID:       r.ID,
		Name:     r.Name,
		Password: pwd,
	})
	if err != nil {
		if errors.Is(err, database.ErrResourceExisted) {
			c.StatusCode(iris.StatusConflict)
			c.WriteString("the user id has existed")
			return
		}
		c.StatusCode(iris.StatusInternalServerError)
		c.WriteString(err.Error())
		return
	}

	c.StatusCode(iris.StatusOK)
}

func (s *Server) UpdateConfig(ctx iris.Context) {
	var r httpModels.UserConfigRequestBody
	if err := ctx.ReadJSON(&r); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	tokenClaims := s.auth.GetTokenClaims(ctx)

	err := s.db.User().UpdateConfig(tokenClaims.UserID, dbModels.UserConfig{
		CompareItemsInDifferentShop: r.CompareItemsInDifferentShop,
		CompareItemsInSameShop:      r.CompareItemsInSameShop,
	})
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	ctx.StatusCode(iris.StatusOK)
}

func (s *Server) GetConfig(ctx iris.Context) {
	tokenClaims := s.auth.GetTokenClaims(ctx)

	cfg, err := s.db.User().GetConfig(tokenClaims.UserID)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	ctx.StatusCode(iris.StatusOK)
	if err := ctx.JSON(cfg); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
	}
}
