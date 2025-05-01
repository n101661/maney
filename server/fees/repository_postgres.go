package fees

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

func NewPostgresRepository(engine *xorm.Engine) (repository.FeeRepository, error) {
	return &postgresRepository{
		engine: engine,
	}, nil
}

func (repo *postgresRepository) Create(ctx context.Context, r *repository.CreateFeesRequest) ([]*repository.Fee, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	rows := lo.Map(r.Fees, func(item *repository.BaseCreateFee, _ int) *postgres.FeesModel {
		return &postgres.FeesModel{
			PublicID: item.PublicID,
			UserID:   r.UserID,
			Name:     item.Name,
			Data:     toPostgresBaseFee(item.BaseFee),
		}
	})
	_, err := session.Insert(rows)
	if err != nil {
		if postgres.UniqueViolationError(err) {
			return nil, repository.ErrDataExists
		}
		return nil, err
	}
	return lo.Map(rows, func(item *postgres.FeesModel, _ int) *repository.Fee {
		return toRepositoryFee(item)
	}), nil
}

func (repo *postgresRepository) List(ctx context.Context, r *repository.ListFeesRequest) (*repository.ListFeesReply, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	var rows []*postgres.FeesModel
	err := session.Find(&rows, &postgres.FeesModel{
		UserID: r.UserID,
	})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, repository.ErrDataNotFound
	}
	return &repository.ListFeesReply{
		Fees: lo.Map(rows, func(item *postgres.FeesModel, _ int) *repository.Fee {
			return toRepositoryFee(item)
		}),
	}, nil
}

func (repo *postgresRepository) Update(ctx context.Context, r *repository.UpdateFeeRequest) (*repository.Fee, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	row := postgres.FeesModel{
		PublicID: r.FeePublicID,
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
	bean := postgres.FeesModel{}
	if r.Fee != nil {
		cols = append(cols, "name", "data")

		row.Name = r.Fee.Name
		bean.Name = row.Name

		row.Data = toPostgresBaseFee(r.Fee)
		bean.Data = row.Data
	}

	affected, err := session.Cols(cols...).Update(&bean, &postgres.FeesModel{
		ID: row.ID,
	})
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, repository.ErrDataNotFound
	}

	return toRepositoryFee(&row), nil
}

func (repo *postgresRepository) Delete(ctx context.Context, r *repository.DeleteFeesRequest) ([]*repository.Fee, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	session.Where("user_id = ?", r.UserID)
	if len(r.FeePublicIDs) > 0 {
		session.In("public_id", r.FeePublicIDs)
	}

	var rows []*postgres.FeesModel
	err := session.Find(&rows)
	if err != nil {
		return nil, err
	}

	if len(r.FeePublicIDs) > 0 && len(rows) != len(r.FeePublicIDs) {
		return nil, repository.ErrDataNotFound
	}
	if len(r.FeePublicIDs) == 0 && len(rows) == 0 {
		return nil, repository.ErrDataNotFound
	}

	_, err = session.In("id", lo.Map(rows, func(item *postgres.FeesModel, _ int) any {
		return item.ID
	})).Table(&postgres.FeesModel{}).Delete()
	if err != nil {
		return nil, err
	}

	return lo.Map(rows, func(item *postgres.FeesModel, _ int) *repository.Fee {
		return toRepositoryFee(item)
	}), nil
}

func toPostgresBaseFee(item *repository.BaseFee) *postgres.BaseFee {
	return &postgres.BaseFee{
		Type:  item.Type,
		Rate:  item.Rate,
		Fixed: item.Fixed,
	}
}

func toRepositoryFee(item *postgres.FeesModel) *repository.Fee {
	return &repository.Fee{
		ID:       item.ID,
		PublicID: item.PublicID,
		BaseFee: &repository.BaseFee{
			Name:  item.Name,
			Type:  item.Data.Type,
			Rate:  item.Data.Rate,
			Fixed: item.Data.Fixed,
		},
	}
}
