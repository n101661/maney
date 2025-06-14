package accounts

import (
	"errors"
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"

	irisController "github.com/n101661/maney/server/controller/iris"
	"github.com/n101661/maney/server/models"
)

type IrisController struct {
	s Service
	*irisController.SimpleCreateTemplate[models.BasicAccount, CreateRequest, CreateReply, models.ObjectId]
	*irisController.SimpleListTemplate[ListRequest, ListReply, []*models.Account]
	*irisController.SimpleUpdateTemplate[models.BasicAccount, UpdateRequest, UpdateReply]
	*irisController.SimpleDeleteTemplate[DeleteRequest, DeleteReply]
}

func NewIrisController(s Service) *IrisController {
	return &IrisController{
		s: s,
		SimpleCreateTemplate: &irisController.SimpleCreateTemplate[models.BasicAccount, CreateRequest, CreateReply, models.ObjectId]{
			Service: s,
			ParseServiceRequest: func(userID string, r *models.BasicAccount) (*CreateRequest, error) {
				initialBalance, err := decimal.NewFromString(r.InitialBalance)
				if err != nil {
					return nil, fmt.Errorf("invalid decimal[%s]", r.InitialBalance)
				}
				return &CreateRequest{
					UserID: userID,
					Account: &BaseAccount{
						Name:           r.Name,
						IconID:         int32(r.IconId),
						InitialBalance: initialBalance,
					},
				}, nil
			},
			BadRequest: func(err error) (httpCode int, yes bool) {
				switch {
				case errors.Is(err, ErrDataInsufficient):
					return iris.StatusBadRequest, true
				}
				return 0, false
			},
			ParseAPIResponse: func(reply *CreateReply) (*models.ObjectId, error) {
				return &models.ObjectId{
					Id: lo.ToPtr(models.Id(reply.Account.PublicID)),
				}, nil
			},
		},
		SimpleListTemplate: &irisController.SimpleListTemplate[ListRequest, ListReply, []*models.Account]{
			Service: s,
			ParseServiceRequest: func(c iris.Context, userID string) (*ListRequest, error) {
				return &ListRequest{
					UserID: userID,
				}, nil
			},
			BadRequest: func(err error) (httpCode int, yes bool) {
				switch {
				case errors.Is(err, ErrDataInsufficient):
					return iris.StatusBadRequest, true
				}
				return 0, false
			},
			ParseAPIResponse: func(reply *ListReply) (*[]*models.Account, error) {
				return lo.ToPtr(lo.Map(reply.Accounts, func(item *Account, _ int) *models.Account {
					return &models.Account{
						Id:             lo.ToPtr(models.Id(item.PublicID)),
						Name:           item.Name,
						IconId:         models.IconId(item.IconID),
						InitialBalance: item.InitialBalance.String(),
						Balance:        lo.ToPtr(item.Balance.String()),
					}
				})), nil
			},
		},
		SimpleUpdateTemplate: &irisController.SimpleUpdateTemplate[models.BasicAccount, UpdateRequest, UpdateReply]{
			Placeholder: "accountId",
			Service:     s,
			ParseServiceRequest: func(userID string, publicID string, r *models.BasicAccount) (*UpdateRequest, error) {
				initialBalance, err := decimal.NewFromString(r.InitialBalance)
				if err != nil {
					return nil, fmt.Errorf("invalid decimal[%s]", r.InitialBalance)
				}
				return &UpdateRequest{
					UserID:          userID,
					AccountPublicID: publicID,
					Account: &BaseAccount{
						Name:           r.Name,
						IconID:         int32(r.IconId),
						InitialBalance: initialBalance,
					},
				}, nil
			},
			BadRequest: func(err error) (httpCode int, yes bool) {
				switch {
				case errors.Is(err, ErrAccountNotFound):
					return iris.StatusNotFound, true
				case errors.Is(err, ErrDataInsufficient):
					return iris.StatusBadRequest, true
				}
				return 0, false
			},
		},
		SimpleDeleteTemplate: &irisController.SimpleDeleteTemplate[DeleteRequest, DeleteReply]{
			Placeholder: "accountId",
			Service:     s,
			ParseServiceRequest: func(userID string, publicID string) *DeleteRequest {
				return &DeleteRequest{
					UserID:          userID,
					AccountPublicID: publicID,
				}
			},
			BadRequest: func(err error) (httpCode int, yes bool) {
				switch {
				case errors.Is(err, ErrAccountNotFound):
					return iris.StatusNotFound, true
				case errors.Is(err, ErrDataInsufficient):
					return iris.StatusBadRequest, true
				}
				return 0, false
			},
		},
	}
}
