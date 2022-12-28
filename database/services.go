package database

import "github.com/n101661/maney/database/models"

type UserService interface {
	// Create returns ErrResourceExisted if there is an existed user.
	Create(models.User) error
	Update(models.User) error
	Get(id string) (*models.User, error)
	UpdateConfig(id string, val models.UserConfig) error
	GetConfig(id string) (models.UserConfig, error)
}

type AccountService interface {
	// Create returns ErrResourceExisted if there is an existed account.
	Create(userID string, account models.AssetAccount) (oid uint64, err error)
	Update(userID string, account models.AssetAccount) error
	Delete(userID string, accountOID uint64) error
	List(userID string) ([]models.AssetAccount, error)
}

type CategoryService interface {
	// Create returns ErrResourceExisted if there is an existed category.
	Create(userID string, category models.Category) (oid uint64, err error)
	Update(userID string, category models.Category) error
	Delete(userID string, categoryOID uint64) error
	Get(userID string, categoryOID uint64) (*models.Category, error)
}

type ShopService interface {
	// Create returns ErrResourceExisted if there is an existed shop.
	Create(userID string, shop models.Shop) (oid uint64, err error)
	Update(userID string, shop models.Shop) error
	Delete(userID string, shopOID uint64) error
	Get(userID string, shopOID uint64) (*models.Shop, error)
}

type FeeService interface {
	// Create returns ErrResourceExisted if there is an existed fee.
	Create(userID string, fee models.Fee) (oid uint64, err error)
	Update(userID string, fee models.Fee) error
	Delete(userID string, feeOID uint64) error
	Get(userID string, feeOID uint64) (*models.Fee, error)
}

type DailyItemService interface {
	// Create returns ErrResourceExisted if there is an existed daily item.
	Create(userID string, item models.DailyItem) (oid uint64, err error)
	CreateMultiple(userID string, items []models.DailyItem) error
	Update(userID string, item models.DailyItem) error
	Delete(userID string, itemOID uint64) error
	List(userID string) ([]models.DailyItem, error)
}

type RepeatingItemService interface {
	// Create returns ErrResourceExisted if there is an existed repeating item.
	Create(userID string, item models.RepeatingItem) (oid uint64, err error)
	Update(userID string, item models.RepeatingItem) error
	Delete(userID string, itemOID uint64) error
	List(userID string) ([]models.RepeatingItem, error)
}
