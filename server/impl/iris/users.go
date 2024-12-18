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
	"github.com/n101661/maney/server/models"
)

func (s *Server) Login(c iris.Context) {
	var r models.LoginRequestBody
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
		"refreshToken", refreshToken.ID,
		iris.CookiePath("/auth"),
		iris.CookieExpires(refreshToken.ExpireAfter),
		iris.CookieHTTPOnly(true),
		iris.CookieSameSite(http.SameSiteStrictMode),
	)

	c.StatusCode(iris.StatusOK)
	err = c.JSON(&models.LoginResponse{
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

func (s *Server) Logout(ctx iris.Context) {
	ctx.Header("Set-Cookie", `token=""; Max-Age=0; HttpOnly`)
	ctx.StatusCode(iris.StatusOK)
}

func (s *Server) SignUp(ctx iris.Context) {
	var r models.SignUpRequestBody
	if err := ctx.ReadJSON(&r); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	pwd, err := s.auth.EncryptPassword(r.Password)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	err = s.db.User().Create(dbModels.User{
		ID:       r.ID,
		Name:     r.Name,
		Password: pwd,
	})
	if err != nil {
		if errors.Is(err, database.ErrResourceExisted) {
			ctx.StatusCode(iris.StatusConflict)
			ctx.WriteString("the user id has existed")
			return
		}
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	ctx.StatusCode(iris.StatusOK)
}

func (s *Server) UpdateConfig(ctx iris.Context) {
	var r models.UserConfigRequestBody
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
