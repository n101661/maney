package categories

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

func NewPostgresRepository(engine *xorm.Engine) (repository.CategoryRepository, error) {
	return &postgresRepository{
		engine: engine,
	}, nil
}

func (repo *postgresRepository) Create(ctx context.Context, r *repository.CreateCategoriesRequest) ([]*repository.Category, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	rows := lo.Map(r.Categories, func(item *repository.BaseCreateCategory, _ int) *postgres.CategoriesModel {
		return &postgres.CategoriesModel{
			PublicID: item.PublicID,
			UserID:   r.UserID,
			Type:     r.Type,
			Data: &postgres.BaseCategory{
				BaseCategory: item.BaseCategory,
			},
		}
	})
	_, err := session.Insert(rows)
	if err != nil {
		if postgres.UniqueViolationError(err) {
			return nil, repository.ErrDataExists
		}
		return nil, err
	}
	return lo.Map(rows, func(item *postgres.CategoriesModel, _ int) *repository.Category {
		return toCategory(item)
	}), nil
}

func (repo *postgresRepository) List(ctx context.Context, r *repository.ListCategoriesRequest) (*repository.ListCategoriesReply, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	var rows []*postgres.CategoriesModel
	err := session.Find(&rows, &postgres.CategoriesModel{
		UserID: r.UserID,
		Type:   r.Type,
	})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, repository.ErrDataNotFound
	}
	return &repository.ListCategoriesReply{
		Categories: lo.Map(rows, func(item *postgres.CategoriesModel, _ int) *repository.Category {
			return toCategory(item)
		}),
	}, nil
}

func (repo *postgresRepository) Update(ctx context.Context, r *repository.UpdateCategoryRequest) (*repository.Category, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	row := postgres.CategoriesModel{
		PublicID: r.CategoryPublicID,
		UserID:   r.UserID,
	}
	has, err := session.Get(&row)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, repository.ErrDataNotFound
	}
	row.Data.BaseCategory = r.Category

	affected, err := session.Update(&postgres.CategoriesModel{
		Data: row.Data,
	}, &postgres.CategoriesModel{
		ID: row.ID,
	})
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, repository.ErrDataNotFound
	}

	return toCategory(&row), nil
}

func (repo *postgresRepository) Delete(ctx context.Context, r *repository.DeleteCategoriesRequest) ([]*repository.Category, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	session.Where("user_id = ?", r.UserID)
	if len(r.CategoryPublicIDs) > 0 {
		session.In("public_id", r.CategoryPublicIDs)
	}

	var rows []*postgres.CategoriesModel
	err := session.Find(&rows)
	if err != nil {
		return nil, err
	}

	if len(r.CategoryPublicIDs) > 0 && len(rows) != len(r.CategoryPublicIDs) {
		return nil, repository.ErrDataNotFound
	}
	if len(r.CategoryPublicIDs) == 0 && len(rows) == 0 {
		return nil, repository.ErrDataNotFound
	}

	_, err = session.In("id", lo.Map(rows, func(item *postgres.CategoriesModel, _ int) any {
		return item.ID
	})).Table(&postgres.CategoriesModel{}).Delete()
	if err != nil {
		return nil, err
	}

	return lo.Map(rows, func(item *postgres.CategoriesModel, _ int) *repository.Category {
		return toCategory(item)
	}), nil
}

func toCategory(item *postgres.CategoriesModel) *repository.Category {
	return &repository.Category{
		ID:           item.ID,
		PublicID:     item.PublicID,
		BaseCategory: item.Data.BaseCategory,
	}
}
