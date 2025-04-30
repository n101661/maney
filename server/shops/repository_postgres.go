package shops

import (
	"context"

	"github.com/n101661/maney/server/repository"
	"github.com/n101661/maney/server/repository/postgres"
	"github.com/samber/lo"
	"xorm.io/xorm"
)

type postgresRepository struct {
	engine *xorm.Engine
}

func NewPostgresRepository(engine *xorm.Engine) (repository.ShopRepository, error) {
	return &postgresRepository{
		engine: engine,
	}, nil
}

func (repo *postgresRepository) Create(ctx context.Context, r *repository.CreateShopsRequest) ([]*repository.Shop, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	rows := lo.Map(r.Shops, func(item *repository.BaseCreateShop, _ int) *postgres.ShopsModel {
		return &postgres.ShopsModel{
			PublicID: item.PublicID,
			UserID:   r.UserID,
			Name:     item.Name,
			Address:  item.Address,
		}
	})
	_, err := session.Insert(rows)
	if err != nil {
		if postgres.UniqueViolationError(err) {
			return nil, repository.ErrDataExists
		}
		return nil, err
	}
	return lo.Map(rows, func(item *postgres.ShopsModel, _ int) *repository.Shop {
		return toShop(item)
	}), nil
}

func (repo *postgresRepository) List(ctx context.Context, r *repository.ListShopsRequest) (*repository.ListShopsReply, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	var rows []*postgres.ShopsModel
	err := session.Find(&rows, &postgres.ShopsModel{
		UserID: r.UserID,
	})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, repository.ErrDataNotFound
	}
	return &repository.ListShopsReply{
		Shops: lo.Map(rows, func(item *postgres.ShopsModel, _ int) *repository.Shop {
			return toShop(item)
		}),
	}, nil
}

func (repo *postgresRepository) Update(ctx context.Context, r *repository.UpdateShopRequest) (*repository.Shop, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	row := postgres.ShopsModel{
		PublicID: r.ShopPublicID,
		UserID:   r.UserID,
	}
	has, err := session.Get(&row)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, repository.ErrDataNotFound
	}

	cols := []string{}
	bean := postgres.ShopsModel{}
	if r.Shop != nil {
		cols = append(cols, "name", "address")

		row.Name = r.Shop.Name
		bean.Name = row.Name

		row.Address = r.Shop.Address
		bean.Address = row.Address
	}

	affected, err := session.Cols(cols...).Update(&bean, &postgres.ShopsModel{
		ID: row.ID,
	})
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, repository.ErrDataNotFound
	}

	return toShop(&row), nil
}

func (repo *postgresRepository) Delete(ctx context.Context, r *repository.DeleteShopsRequest) ([]*repository.Shop, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	session.Where("user_id = ?", r.UserID)
	if len(r.ShopPublicIDs) > 0 {
		session.In("public_id", r.ShopPublicIDs)
	}

	var rows []*postgres.ShopsModel
	err := session.Find(&rows)
	if err != nil {
		return nil, err
	}

	if len(r.ShopPublicIDs) > 0 && len(rows) != len(r.ShopPublicIDs) {
		return nil, repository.ErrDataNotFound
	}
	if len(r.ShopPublicIDs) == 0 && len(rows) == 0 {
		return nil, repository.ErrDataNotFound
	}

	_, err = session.In("id", lo.Map(rows, func(item *postgres.ShopsModel, _ int) any {
		return item.ID
	})).Table(&postgres.ShopsModel{}).Delete()
	if err != nil {
		return nil, err
	}

	return lo.Map(rows, func(item *postgres.ShopsModel, _ int) *repository.Shop {
		return toShop(item)
	}), nil
}

func toShop(item *postgres.ShopsModel) *repository.Shop {
	return &repository.Shop{
		ID:       item.ID,
		PublicID: item.PublicID,
		BaseShop: &repository.BaseShop{
			Name:    item.Name,
			Address: item.Address,
		},
	}
}
