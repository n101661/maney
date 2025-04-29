package users

import (
	"context"
	"database/sql"
	"time"

	"xorm.io/xorm"

	"github.com/n101661/maney/server/repository"
	"github.com/n101661/maney/server/repository/postgres"
	"github.com/samber/lo"
)

type postgresRepository struct {
	engine *xorm.Engine
}

func NewPostgresRepository(engine *xorm.Engine) (repository.UserRepository, error) {
	return &postgresRepository{
		engine: engine,
	}, nil
}

func (repo *postgresRepository) CreateUser(ctx context.Context, user *repository.UserModel) error {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	_, err := session.Insert(postgres.UsersModel{
		ID:       user.ID,
		Password: user.Password,
		Config: &postgres.UserConfig{
			UserConfig: user.Config,
		},
	})
	if err != nil {
		if postgres.UniqueViolationError(err) {
			return repository.ErrDataExists
		}
		return err
	}
	return nil
}

func (repo *postgresRepository) GetUser(ctx context.Context, userID string) (*repository.UserModel, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	user := postgres.UsersModel{
		ID: userID,
	}
	has, err := session.Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, repository.ErrDataNotFound
	}

	return &repository.UserModel{
		ID:       user.ID,
		Password: user.Password,
		Config:   user.Config.UserConfig,
	}, nil
}

func (repo *postgresRepository) UpdateUser(ctx context.Context, user *repository.UserModel) error {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	effectedRows, err := session.Update(
		postgres.UsersModel{
			Password: user.Password,
			Config: &postgres.UserConfig{
				UserConfig: user.Config,
			},
		},
		postgres.UsersModel{
			ID: user.ID,
		},
	)
	if err != nil {
		return err
	}
	if effectedRows == 0 {
		return repository.ErrDataNotFound
	}
	return nil
}

func (repo *postgresRepository) CreateToken(ctx context.Context, token *repository.TokenModel) error {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	_, err := session.Insert(postgres.TokensModel{
		ID:         token.ID,
		UserID:     token.Claim.UserID,
		ExpiryTime: token.ExpiryTime,
	})
	if err != nil {
		if postgres.UniqueViolationError(err) {
			return repository.ErrDataExists
		}
		return err
	}
	return nil
}

func (repo *postgresRepository) GetToken(ctx context.Context, tokenID string) (*repository.TokenModel, error) {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	token := postgres.TokensModel{
		ID: tokenID,
	}
	has, err := session.Get(&token)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, repository.ErrDataNotFound
	}
	return &repository.TokenModel{
		ID: token.ID,
		Claim: &repository.TokenClaims{
			UserID: token.UserID,
		},
		ExpiryTime: token.ExpiryTime,
		RevokedAt: lo.IfF(token.RevokedAt.Valid, func() *time.Time {
			return &token.RevokedAt.Time
		}).Else(nil),
	}, nil
}

func (repo *postgresRepository) RevokeToken(ctx context.Context, tokenID string) error {
	session := repo.engine.NewSession().Context(ctx)
	defer session.Close()

	effectedRows, err := session.Update(
		postgres.TokensModel{
			RevokedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		},
		postgres.TokensModel{
			ID: tokenID,
		},
	)
	if err != nil {
		return err
	}
	if effectedRows == 0 {
		return repository.ErrDataNotFound
	}
	return nil
}
