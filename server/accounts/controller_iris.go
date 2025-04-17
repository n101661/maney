package accounts

import (
	"errors"

	"github.com/kataras/iris/v12"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"

	"github.com/n101661/maney/server/models"
)

type IrisController struct {
	s Service
}

func NewIrisController(s Service) *IrisController {
	return &IrisController{
		s: s,
	}
}

func (controller *IrisController) Create(c iris.Context) {
	var r models.BasicAccount
	if err := c.ReadJSON(&r); err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	initialBalance, err := decimal.NewFromString(r.InitialBalance)
	if err != nil {
		c.StopWithText(iris.StatusBadRequest, "invalid decimal")
		return
	}

	user := c.User()
	if user == nil {
		c.StopWithJSON(iris.StatusUnauthorized, &models.EmptyResponse{})
		return
	}

	userID, err := user.GetID()
	if err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	reply, err := controller.s.Create(c.Request().Context(), &CreateRequest{
		UserID: userID,
		Account: &BaseAccount{
			Name:           r.Name,
			IconID:         int32(r.IconId),
			InitialBalance: initialBalance,
		},
	})
	if err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	c.StopWithJSON(iris.StatusOK, &models.ObjectId{
		Id: lo.ToPtr(models.Id(reply.Account.ID)),
	})
}

func (controller *IrisController) List(c iris.Context) {
	user := c.User()
	if user == nil {
		c.StopWithJSON(iris.StatusUnauthorized, &models.EmptyResponse{})
		return
	}

	userID, err := user.GetID()
	if err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	reply, err := controller.s.List(c.Request().Context(), &ListRequest{
		UserID: userID,
	})
	if err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	c.StopWithJSON(iris.StatusOK, lo.Map(reply.Accounts, func(item *Account, _ int) *models.Account {
		return &models.Account{
			Id:             lo.ToPtr(models.Id(item.ID)),
			Name:           item.Name,
			IconId:         models.Id(item.IconID),
			InitialBalance: item.InitialBalance.String(),
			Balance:        lo.ToPtr(item.Balance.String()),
		}
	}))
}

func (controller *IrisController) Update(c iris.Context) {
	accountID, err := c.Params().GetInt32("accountId")
	if err != nil {
		c.StopWithText(iris.StatusBadRequest, "missing or invalid account id")
		return
	}

	var r models.BasicAccount
	if err := c.ReadJSON(&r); err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	initialBalance, err := decimal.NewFromString(r.InitialBalance)
	if err != nil {
		c.StopWithText(iris.StatusBadRequest, "invalid decimal")
		return
	}

	user := c.User()
	if user == nil {
		c.StopWithJSON(iris.StatusUnauthorized, &models.EmptyResponse{})
		return
	}

	userID, err := user.GetID()
	if err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	_, err = controller.s.Update(c.Request().Context(), &UpdateRequest{
		UserID:    userID,
		AccountID: accountID,
		Account: &BaseAccount{
			Name:           r.Name,
			IconID:         int32(r.IconId),
			InitialBalance: initialBalance,
		},
	})
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			c.StopWithText(iris.StatusNotFound, "no [%d] account id", accountID)
		} else {
			c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		}
		return
	}

	c.StopWithJSON(iris.StatusOK, &models.EmptyResponse{})
}

func (controller *IrisController) Delete(c iris.Context) {
	accountID, err := c.Params().GetInt32("accountId")
	if err != nil {
		c.StopWithText(iris.StatusBadRequest, "missing or invalid account id")
		return
	}

	user := c.User()
	if user == nil {
		c.StopWithJSON(iris.StatusUnauthorized, &models.EmptyResponse{})
		return
	}

	userID, err := user.GetID()
	if err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	_, err = controller.s.Delete(c.Request().Context(), &DeleteRequest{
		UserID:    userID,
		AccountID: accountID,
	})
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			c.StopWithText(iris.StatusNotFound, "no [%d] account id", accountID)
		} else {
			c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		}
		return
	}

	c.StopWithJSON(iris.StatusOK, &models.EmptyResponse{})
}
