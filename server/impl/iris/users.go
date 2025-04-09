package iris

import (
	"github.com/kataras/iris/v12"

	dbModels "github.com/n101661/maney/database/models"
	httpModels "github.com/n101661/maney/server/models"
)

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
