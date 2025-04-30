package categories

import (
	"errors"

	"github.com/kataras/iris/v12"
	"github.com/samber/lo"

	irisController "github.com/n101661/maney/server/controller/iris"
	"github.com/n101661/maney/server/models"
	"github.com/n101661/maney/server/repository"
)

type IrisController struct {
	s Service
	*irisController.SimpleCreateTemplate[models.CreatingCategory, CreateRequest, CreateReply, models.ObjectId]
	*irisController.SimpleListTemplate[ListRequest, ListReply, []*models.Category]
	*irisController.SimpleUpdateTemplate[models.BasicCategory, UpdateRequest, UpdateReply]
	*irisController.SimpleDeleteTemplate[DeleteRequest, DeleteReply]
}

func NewIrisController(s Service) *IrisController {
	return &IrisController{
		s: s,
		SimpleCreateTemplate: &irisController.SimpleCreateTemplate[models.CreatingCategory, CreateRequest, CreateReply, models.ObjectId]{
			Service: s,
			ParseServiceRequest: func(userID string, r *models.CreatingCategory) (*CreateRequest, error) {
				type_, err := parseType(string(r.Type))
				if err != nil {
					return nil, err
				}
				return &CreateRequest{
					UserID: userID,
					Type:   type_,
					Category: &BaseCategory{
						Name:   r.Name,
						IconID: int32(lo.FromPtrOr(r.IconId, 0)),
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
					Id: lo.ToPtr(models.Id(reply.Category.PublicID)),
				}, nil
			},
		},
		SimpleListTemplate: &irisController.SimpleListTemplate[ListRequest, ListReply, []*models.Category]{
			Service: s,
			ParseServiceRequest: func(c iris.Context, userID string) (*ListRequest, error) {
				type_, err := parseType(c.URLParamDefault("type", repository.CategoryTypeExpense.String()))
				if err != nil {
					return nil, err
				}
				return &ListRequest{
					UserID: userID,
					Type:   type_,
				}, nil
			},
			BadRequest: func(err error) (httpCode int, yes bool) {
				switch {
				case errors.Is(err, ErrDataInsufficient):
					return iris.StatusBadRequest, true
				}
				return 0, false
			},
			ParseAPIResponse: func(reply *ListReply) (*[]*models.Category, error) {
				return lo.ToPtr(lo.Map(reply.Categories, func(item *Category, _ int) *models.Category {
					return &models.Category{
						Id:     lo.ToPtr(models.Id(item.PublicID)),
						Name:   item.Name,
						IconId: lo.ToPtr(models.IconId(item.IconID)),
					}
				})), nil
			},
		},
		SimpleUpdateTemplate: &irisController.SimpleUpdateTemplate[models.BasicCategory, UpdateRequest, UpdateReply]{
			Placeholder: "categoryId",
			Service:     s,
			ParseServiceRequest: func(userID string, publicID string, r *models.BasicCategory) (*UpdateRequest, error) {
				return &UpdateRequest{
					UserID:           userID,
					CategoryPublicID: publicID,
					Category: &BaseCategory{
						Name:   r.Name,
						IconID: int32(lo.FromPtrOr(r.IconId, 0)),
					},
				}, nil
			},
			BadRequest: func(err error) (httpCode int, yes bool) {
				switch {
				case errors.Is(err, ErrCategoryNotFound):
					return iris.StatusNotFound, true
				case errors.Is(err, ErrDataInsufficient):
					return iris.StatusBadRequest, true
				}
				return 0, false
			},
		},
		SimpleDeleteTemplate: &irisController.SimpleDeleteTemplate[DeleteRequest, DeleteReply]{
			Placeholder: "categoryId",
			Service:     s,
			ParseServiceRequest: func(userID string, publicID string) *DeleteRequest {
				return &DeleteRequest{
					UserID:           userID,
					CategoryPublicID: publicID,
				}
			},
			BadRequest: func(err error) (httpCode int, yes bool) {
				switch {
				case errors.Is(err, ErrCategoryNotFound):
					return iris.StatusNotFound, true
				case errors.Is(err, ErrDataInsufficient):
					return iris.StatusBadRequest, true
				}
				return 0, false
			},
		},
	}
}

func parseType(s string) (Type, error) {
	return repository.ToCategoryType(s)
}
