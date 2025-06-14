package postgres

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"

	"github.com/n101661/maney/server/repository"
)

type UsersModel struct {
	ID        string      `xorm:"pk"`
	Password  []byte      `xorm:"not null"`
	Config    *UserConfig `xorm:"json not null"`
	CreatedAt time.Time   `xorm:"created not null"`
}

func (*UsersModel) TableName() string {
	return "users"
}

type UserConfig struct {
	*repository.UserConfig
}

func (v *UserConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.UserConfig)
}

func (v *UserConfig) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.UserConfig)
}

type TokensModel struct {
	ID         string       `xorm:"char(88) pk"`
	UserID     string       `xorm:"index not null"`
	ExpiryTime time.Time    `xorm:"not null"`
	CreatedAt  time.Time    `xorm:"created not null"`
	RevokedAt  sql.NullTime `xorm:"timestamp null"`
}

func (*TokensModel) TableName() string {
	return "tokens"
}

type AccountsModel struct {
	ID       int32               `xorm:"serial pk"`
	PublicID string              `xorm:"unique not null"`
	UserID   string              `xorm:"index not null"`
	Data     *BaseAccount        `xorm:"json not null"`
	Balance  decimal.NullDecimal `xorm:"numeric(15,6) not null"`
}

func (*AccountsModel) TableName() string {
	return "accounts"
}

type BaseAccount struct {
	*repository.BaseAccount
}

func (v *BaseAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.BaseAccount)
}

func (v *BaseAccount) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.BaseAccount)
}

type CategoriesModel struct {
	ID       int32                   `xorm:"serial pk"`
	PublicID string                  `xorm:"unique not null"`
	UserID   string                  `xorm:"index not null"`
	Type     repository.CategoryType `xorm:"smallint not null"`
	Data     *BaseCategory           `xorm:"json not null"`
}

func (*CategoriesModel) TableName() string {
	return "categories"
}

type BaseCategory struct {
	*repository.BaseCategory
}

func (v *BaseCategory) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.BaseCategory)
}

func (v *BaseCategory) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.BaseCategory)
}

type ShopsModel struct {
	ID       int32  `xorm:"serial pk"`
	PublicID string `xorm:"unique not null"`
	UserID   string `xorm:"index not null"`
	Name     string `xorm:"text not null"`
	Address  string `xorm:"text not null"`
}

func (*ShopsModel) TableName() string {
	return "shops"
}

type FeesModel struct {
	ID       int32    `xorm:"serial pk"`
	PublicID string   `xorm:"unique not null"`
	UserID   string   `xorm:"index not null"`
	Name     string   `xorm:"text not null"`
	Data     *BaseFee `xorm:"json not null"`
}

func (*FeesModel) TableName() string {
	return "fees"
}

type BaseFee struct {
	Type  int8
	Rate  *decimal.Decimal `json:"rate,omitempty"`
	Fixed *decimal.Decimal `json:"fixed,omitempty"`
}

func (v *BaseFee) MarshalJSON() ([]byte, error) {
	t := tempBaseFee(*v)
	return json.Marshal(t)
}

func (v *BaseFee) UnmarshalJSON(data []byte) error {
	var t tempBaseFee
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	*v = BaseFee(t)
	return nil
}

type tempBaseFee BaseFee
