package database

import "github.com/n101661/maney/database/models"

type UserService interface {
	// Create returns ErrResourceExisted if there is an existed user.
	Create(models.User) error
	Update(models.User) error
	Get(id string) (*models.User, error)
}

type AccountService interface {
	// Create returns ErrResourceExisted if there is an existed account.
	Create(userOID uint64, account models.AssetAccount) error
	Update(userOID uint64, account models.AssetAccount) error
	Delete(userOID, accountOID uint64) error
	Get(userOID, accountOID uint64) (*models.AssetAccount, error)
}

type CategoryService interface {
	// Create returns ErrResourceExisted if there is an existed category.
	Create(userOID uint64, category models.Category) error
	Update(userOID uint64, category models.Category) error
	Delete(userOID, categoryOID uint64) error
	Get(userOID, categoryOID uint64) (*models.Category, error)
}

type ShopService interface {
	// Create returns ErrResourceExisted if there is an existed shop.
	Create(userOID uint64, shop models.Shop) error
	Update(userOID uint64, shop models.Shop) error
	Delete(userOID, shopOID uint64) error
	Get(userOID, shopOID uint64) (*models.Shop, error)
}

type FeeService interface {
	// Create returns ErrResourceExisted if there is an existed fee.
	Create(userOID uint64, fee models.Fee) error
	Update(userOID uint64, fee models.Fee) error
	Delete(userOID, feeOID uint64) error
	Get(userOID, feeOID uint64) (*models.Fee, error)
}

type DailyItemService interface {
	// Create returns ErrResourceExisted if there is an existed daily item.
	Create(userOID uint64, item models.DailyItem) error
	CreateMultiple(userOID uint64, items []models.DailyItem) error
	Update(userOID uint64, item models.DailyItem) error
	Delete(userOID, itemOID uint64) error
	List(userOID uint64) ([]models.DailyItem, error)
}

type RepeatingItemService interface {
	// Create returns ErrResourceExisted if there is an existed repeating item.
	Create(userOID uint64, item models.RepeatingItem) error
	Update(userOID uint64, item models.RepeatingItem) error
	Delete(userOID, itemOID uint64) error
	List(userOID uint64) ([]models.RepeatingItem, error)
}
