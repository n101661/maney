package shops

import (
	"errors"

	"github.com/kataras/iris/v12"
	"github.com/samber/lo"

	irisController "github.com/n101661/maney/server/controller/iris"
	"github.com/n101661/maney/server/models"
)

type IrisController struct {
	*irisController.SimpleCreateTemplate[models.BasicShop, CreateRequest, CreateReply, models.ObjectId]
	*irisController.SimpleListTemplate[ListRequest, ListReply, []*models.Shop]
	*irisController.SimpleUpdateTemplate[models.BasicShop, UpdateRequest, UpdateReply]
	*irisController.SimpleDeleteTemplate[DeleteRequest, DeleteReply]
}

func NewIrisController(s Service) *IrisController {
	return &IrisController{
		SimpleCreateTemplate: &irisController.SimpleCreateTemplate[models.BasicShop, CreateRequest, CreateReply, models.ObjectId]{
			Service: s,
			ParseServiceRequest: func(userID string, r *models.BasicShop) (*CreateRequest, error) {
				return &CreateRequest{
					UserID: userID,
					Shop: &BaseShop{
						Name:    r.Name,
						Address: lo.FromPtr(r.Address),
					},
				}, nil
			},
			ParseAPIResponse: func(reply *CreateReply) (*models.ObjectId, error) {
				return &models.ObjectId{
					Id: lo.ToPtr(models.Id(reply.Shop.PublicID)),
				}, nil
			},
		},
		SimpleListTemplate: &irisController.SimpleListTemplate[ListRequest, ListReply, []*models.Shop]{
			Service: s,
			ParseServiceRequest: func(c iris.Context, userID string) (*ListRequest, error) {
				return &ListRequest{
					UserID: userID,
				}, nil
			},
			ParseAPIResponse: func(reply *ListReply) (*[]*models.Shop, error) {
				return lo.ToPtr(lo.Map(reply.Shops, func(item *Shop, _ int) *models.Shop {
					return &models.Shop{
						Id:      lo.ToPtr(models.Id(item.PublicID)),
						Name:    item.Name,
						Address: lo.ToPtr(item.Address),
					}
				})), nil
			},
		},
		SimpleUpdateTemplate: &irisController.SimpleUpdateTemplate[models.BasicShop, UpdateRequest, UpdateReply]{
			Placeholder: "shopId",
			Service:     s,
			ParseServiceRequest: func(userID string, publicID string, r *models.BasicShop) (*UpdateRequest, error) {
				return &UpdateRequest{
					UserID:       userID,
					ShopPublicID: publicID,
					Shop: &BaseShop{
						Name:    r.Name,
						Address: lo.FromPtr(r.Address),
					},
				}, nil
			},
			ResourceNotFound: func(err error) bool {
				return errors.Is(err, ErrShopNotFound)
			},
		},
		SimpleDeleteTemplate: &irisController.SimpleDeleteTemplate[DeleteRequest, DeleteReply]{
			Placeholder: "shopId",
			Service:     s,
			ParseServiceRequest: func(userID string, publicID string) *DeleteRequest {
				return &DeleteRequest{
					UserID:       userID,
					ShopPublicID: publicID,
				}
			},
			ResourceNotFound: func(err error) bool {
				return errors.Is(err, ErrShopNotFound)
			},
		},
	}
}
