package iris

import (
	"errors"

	"github.com/kataras/iris/v12"
	"github.com/shopspring/decimal"

	"github.com/n101661/maney/database"
	dbModels "github.com/n101661/maney/database/models"
	"github.com/n101661/maney/server/models"
)

func (s *Server) CreateAccount(ctx iris.Context) {
	var req models.Account
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	token := s.auth.GetTokenClaims(ctx)

	oid, err := s.db.Account().Create(token.UserID, dbModels.AssetAccount{
		Name:           req.Name,
		IconOID:        req.IconOID,
		InitialBalance: req.InitialBalance,
		Balance:        decimal.Zero,
	})
	if err != nil {
		if errors.Is(err, database.ErrResourceExisted) {
			ctx.StopWithError(iris.StatusConflict, err)
			return
		}
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	ctx.StatusCode(iris.StatusOK)

	err = ctx.JSON(models.ObjectOID{
		OID: oid,
	})
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
	}
}

func (s *Server) ListAccounts(ctx iris.Context) {
	token := s.auth.GetTokenClaims(ctx)

	accounts, err := s.db.Account().List(token.UserID)
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	resp := make([]models.GetAccountResponse, len(accounts))
	for i, acc := range accounts {
		resp[i] = models.GetAccountResponse{
			OID: acc.OID,
			Account: models.Account{
				Name:           acc.Name,
				IconOID:        acc.IconOID,
				InitialBalance: acc.InitialBalance,
			},
			Balance: acc.Balance,
		}
	}

	ctx.StatusCode(iris.StatusOK)

	if err = ctx.JSON(resp); err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
	}
}

func (s *Server) UpdateAccount(ctx iris.Context) {
	oid, err := ctx.Params().GetUint64("oid")
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	var body models.Account
	if err := ctx.ReadJSON(&body); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	token := s.auth.GetTokenClaims(ctx)

	err = s.db.Account().Update(token.UserID, dbModels.AssetAccount{
		OID:            oid,
		Name:           body.Name,
		IconOID:        body.IconOID,
		InitialBalance: body.InitialBalance,
	})
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	ctx.StopWithStatus(iris.StatusOK)
}

func (s *Server) DeleteAccount(ctx iris.Context) {
	oid, err := ctx.Params().GetUint64("oid")
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	token := s.auth.GetTokenClaims(ctx)

	if err = s.db.Account().Delete(token.UserID, oid); err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	ctx.StopWithStatus(iris.StatusOK)
}
