package shops

import (
	"context"
	"errors"
	"fmt"

	"github.com/n101661/maney/pkg/utils"
	"github.com/n101661/maney/pkg/utils/slugid"
	"github.com/n101661/maney/server/repository"
	"github.com/samber/lo"
)

type service struct {
	repository repository.ShopRepository

	opts *ShopServiceOptions
}

func NewService(
	repository repository.ShopRepository,
	opts ...utils.Option[ShopServiceOptions],
) (Service, error) {
	return &service{
		repository: repository,
		opts:       utils.ApplyOptions(defaultShopServiceOptions(), opts),
	}, nil
}

func (s *service) Create(ctx context.Context, r *CreateRequest) (*CreateReply, error) {
	if r.UserID == "" {
		return nil, fmt.Errorf("missing user id")
	}

	rows, err := s.repository.Create(ctx, &repository.CreateShopsRequest{
		UserID: r.UserID,
		Shops: []*repository.BaseCreateShop{
			parseBaseCreateShop(r.Shop, s.opts.genPublicID),
		},
	})
	if err != nil {
		return nil, err
	}

	return &CreateReply{
		Shop: parseShop(rows[0]),
	}, nil
}

func (s *service) List(ctx context.Context, r *ListRequest) (*ListReply, error) {
	reply, err := s.repository.List(ctx, &repository.ListShopsRequest{
		UserID: r.UserID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return &ListReply{
				Shops: []*Shop{},
			}, nil
		}
		return nil, err
	}

	return &ListReply{
		Shops: lo.Map(reply.Shops, func(item *repository.Shop, _ int) *Shop {
			return parseShop(item)
		}),
	}, nil
}

func (s *service) Update(ctx context.Context, r *UpdateRequest) (*UpdateReply, error) {
	if r.Shop == nil {
		return nil, fmt.Errorf("nothing to update")
	}

	row, err := s.repository.Update(ctx, &repository.UpdateShopRequest{
		UserID:       r.UserID,
		ShopPublicID: r.ShopPublicID,
		Shop:         parseBaseShop(r.Shop),
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return nil, ErrShopNotFound
		}
		return nil, err
	}

	return &UpdateReply{
		Shop: parseShop(row),
	}, nil
}

func (s *service) Delete(ctx context.Context, r *DeleteRequest) (*DeleteReply, error) {
	_, err := s.repository.Delete(ctx, &repository.DeleteShopsRequest{
		ShopPublicIDs: []string{r.ShopPublicID},
		UserID:        r.UserID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return nil, ErrShopNotFound
		}
		return nil, err
	}
	return &DeleteReply{}, nil
}

func parseShop(v *repository.Shop) *Shop {
	return &Shop{
		ID:       v.ID,
		PublicID: v.PublicID,
		BaseShop: &BaseShop{
			Name:    v.Name,
			Address: v.Address,
		},
	}
}

func parseBaseShop(v *BaseShop) *repository.BaseShop {
	if v == nil {
		return nil
	}
	return &repository.BaseShop{
		Name:    v.Name,
		Address: v.Address,
	}
}

func parseBaseCreateShop(v *BaseShop, genPublicID func() string) *repository.BaseCreateShop {
	if v == nil {
		return nil
	}
	return &repository.BaseCreateShop{
		PublicID: genPublicID(),
		BaseShop: parseBaseShop(v),
	}
}

type ShopServiceOptions struct {
	genPublicID func() string
}

func defaultShopServiceOptions() *ShopServiceOptions {
	return &ShopServiceOptions{
		genPublicID: func() string {
			return slugid.New("shp", 11)
		},
	}
}

func WithShopServiceGenPublicID(f func() string) utils.Option[ShopServiceOptions] {
	return func(o *ShopServiceOptions) {
		o.genPublicID = f
	}
}
