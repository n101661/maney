package iris

import (
	"errors"
	"fmt"
	"time"

	"github.com/kataras/iris/v12"

	"github.com/n101661/maney/database"
	dbModels "github.com/n101661/maney/database/models"
	"github.com/n101661/maney/server/impl/iris/auth"
	"github.com/n101661/maney/server/models"
)

func (s *Server) Login(ctx iris.Context) {
	var r models.LogInRequestBody
	if err := ctx.ReadJSON(&r); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	user, err := s.db.User().Get(r.ID)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	if user == nil ||
		s.auth.ValidatePassword(user.Password, []byte(r.Password)) != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}

	tokenMaxAge := 8 * time.Hour

	token, err := s.auth.GenerateToken(auth.TokenClaims{
		UserID: user.ID,
		Name:   user.Name,
	}, time.Now().Add(tokenMaxAge))
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	ctx.Header(
		"Set-Cookie",
		fmt.Sprintf("token=%s; Max-Age=%d; HttpOnly", token, tokenMaxAge/time.Second),
	)
	ctx.StatusCode(iris.StatusOK)
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
