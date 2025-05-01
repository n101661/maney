package fees

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
	*irisController.SimpleCreateTemplate[models.BasicFee, CreateRequest, CreateReply, models.ObjectId]
	*irisController.SimpleListTemplate[ListRequest, ListReply, []*models.Fee]
	*irisController.SimpleUpdateTemplate[models.BasicFee, UpdateRequest, UpdateReply]
	*irisController.SimpleDeleteTemplate[DeleteRequest, DeleteReply]
}

func NewIrisController(s Service) *IrisController {
	return &IrisController{
		SimpleCreateTemplate: &irisController.SimpleCreateTemplate[models.BasicFee, CreateRequest, CreateReply, models.ObjectId]{
			Service: s,
			ParseServiceRequest: func(userID string, r *models.BasicFee) (*CreateRequest, error) {
				base, err := toServiceBaseFee(r)
				if err != nil {
					return nil, err
				}
				return &CreateRequest{
					UserID: userID,
					Fee:    base,
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
					Id: lo.ToPtr(models.Id(reply.Fee.PublicID)),
				}, nil
			},
		},
		SimpleListTemplate: &irisController.SimpleListTemplate[ListRequest, ListReply, []*models.Fee]{
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
			ParseAPIResponse: func(reply *ListReply) (*[]*models.Fee, error) {
				result := make([]*models.Fee, len(reply.Fees))
				for i, item := range reply.Fees {
					f, err := toFee(item)
					if err != nil {
						return nil, err
					}
					result[i] = f
				}
				return &result, nil
			},
		},
		SimpleUpdateTemplate: &irisController.SimpleUpdateTemplate[models.BasicFee, UpdateRequest, UpdateReply]{
			Placeholder: "feeId",
			Service:     s,
			ParseServiceRequest: func(userID string, publicID string, r *models.BasicFee) (*UpdateRequest, error) {
				f, err := toServiceBaseFee(r)
				if err != nil {
					return nil, err
				}
				return &UpdateRequest{
					UserID:      userID,
					FeePublicID: publicID,
					Fee:         f,
				}, nil
			},
			BadRequest: func(err error) (httpCode int, yes bool) {
				switch {
				case errors.Is(err, ErrFeeNotFound):
					return iris.StatusNotFound, true
				case errors.Is(err, ErrDataInsufficient):
					return iris.StatusBadRequest, true
				}
				return 0, false
			},
		},
		SimpleDeleteTemplate: &irisController.SimpleDeleteTemplate[DeleteRequest, DeleteReply]{
			Placeholder: "feeId",
			Service:     s,
			ParseServiceRequest: func(userID string, publicID string) *DeleteRequest {
				return &DeleteRequest{
					UserID:      userID,
					FeePublicID: publicID,
				}
			},
			BadRequest: func(err error) (httpCode int, yes bool) {
				switch {
				case errors.Is(err, ErrFeeNotFound):
					return iris.StatusNotFound, true
				case errors.Is(err, ErrDataInsufficient):
					return iris.StatusBadRequest, true
				}
				return 0, false
			},
		},
	}
}

func toServiceBaseFee(base *models.BasicFee) (*BaseFee, error) {
	result := &BaseFee{
		Name: base.Name,
		Type: int8(base.Type),
	}

	switch base.Type {
	case models.BasicFeeTypeRate:
		v, err := base.Value.AsBasicFeeValue0()
		if err != nil {
			return nil, irisController.InternalError(err)
		}

		if v.Rate != nil {
			rate, err := decimal.NewFromString(*v.Rate)
			if err != nil {
				return nil, fmt.Errorf("invalid decimal[%s]", *v.Rate)
			}
			result.Rate = &rate
		}
	case models.BasicFeeTypeFixed:
		v, err := base.Value.AsBasicFeeValue1()
		if err != nil {
			return nil, irisController.InternalError(err)
		}

		if v.Fixed != nil {
			fixed, err := decimal.NewFromString(*v.Fixed)
			if err != nil {
				return nil, fmt.Errorf("invalid decimal[%s]", *v.Fixed)
			}
			result.Fixed = &fixed
		}
	default:
		return nil, fmt.Errorf("unknown fee type")
	}
	return result, nil
}

func toFee(v *Fee) (*models.Fee, error) {
	result := &models.Fee{
		Id:   lo.ToPtr(models.Id(v.PublicID)),
		Name: v.Name,
		Type: models.FeeType(v.Type),
	}

	switch result.Type {
	case models.FeeTypeRate:
		err := result.Value.FromFeeValue0(models.FeeValue0{
			Rate: lo.ToPtr(models.Decimal(v.BaseFee.Rate.String())),
		})
		if err != nil {
			return nil, err
		}
	case models.FeeTypeFixed:
		err := result.Value.FromFeeValue1(models.FeeValue1{
			Fixed: lo.ToPtr(models.Decimal(v.BaseFee.Fixed.String())),
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown fee type[%d]", result.Type)
	}
	return result, nil
}
