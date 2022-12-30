package iris

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"

	"github.com/n101661/maney/database"
	dbModels "github.com/n101661/maney/database/models"
)

type mockAccountService struct{ mock.Mock }

// Create implements database.AccountService
func (s *mockAccountService) Create(userID string, account dbModels.AssetAccount) (uint64, error) {
	if account.InitialBalance.IsZero() {
		account.InitialBalance = decimal.Zero
	}

	args := s.Called(userID, account)
	return args.Get(0).(uint64), args.Error(1)
}

// Delete implements database.AccountService
func (s *mockAccountService) Delete(userID string, accountOID uint64) error {
	args := s.Called(userID, accountOID)
	return args.Error(0)
}

// List implements database.AccountService
func (s *mockAccountService) List(userID string) ([]dbModels.AssetAccount, error) {
	args := s.Called(userID)
	return args.Get(0).([]dbModels.AssetAccount), args.Error(1)
}

// Update implements database.AccountService
func (s *mockAccountService) Update(userID string, account dbModels.AssetAccount) error {
	args := s.Called(userID, account)
	return args.Error(0)
}

type mockCategoryService struct{ mock.Mock }

// Create implements database.CategoryService
func (*mockCategoryService) Create(userID string, category dbModels.Category) (uint64, error) {
	panic("unimplemented")
}

// Delete implements database.CategoryService
func (*mockCategoryService) Delete(userID string, categoryOID uint64) error {
	panic("unimplemented")
}

// Get implements database.CategoryService
func (*mockCategoryService) Get(userID string, categoryOID uint64) (*dbModels.Category, error) {
	panic("unimplemented")
}

// Update implements database.CategoryService
func (*mockCategoryService) Update(userID string, category dbModels.Category) error {
	panic("unimplemented")
}

type mockDailyItemService struct{ mock.Mock }

// Create implements database.DailyItemService
func (*mockDailyItemService) Create(userID string, item dbModels.DailyItem) (uint64, error) {
	panic("unimplemented")
}

// CreateMultiple implements database.DailyItemService
func (*mockDailyItemService) CreateMultiple(userID string, items []dbModels.DailyItem) error {
	panic("unimplemented")
}

// Delete implements database.DailyItemService
func (*mockDailyItemService) Delete(userID string, itemOID uint64) error {
	panic("unimplemented")
}

// List implements database.DailyItemService
func (*mockDailyItemService) List(userID string) ([]dbModels.DailyItem, error) {
	panic("unimplemented")
}

// Update implements database.DailyItemService
func (*mockDailyItemService) Update(userID string, item dbModels.DailyItem) error {
	panic("unimplemented")
}

type mockFeeService struct{ mock.Mock }

// Create implements database.FeeService
func (*mockFeeService) Create(userID string, fee dbModels.Fee) (uint64, error) {
	panic("unimplemented")
}

// Delete implements database.FeeService
func (*mockFeeService) Delete(userID string, feeOID uint64) error {
	panic("unimplemented")
}

// Get implements database.FeeService
func (*mockFeeService) Get(userID string, feeOID uint64) (*dbModels.Fee, error) {
	panic("unimplemented")
}

// Update implements database.FeeService
func (*mockFeeService) Update(userID string, fee dbModels.Fee) error {
	panic("unimplemented")
}

type mockRepeatingItemService struct{ mock.Mock }

// Create implements database.RepeatingItemService
func (*mockRepeatingItemService) Create(userID string, item dbModels.RepeatingItem) (uint64, error) {
	panic("unimplemented")
}

// Delete implements database.RepeatingItemService
func (*mockRepeatingItemService) Delete(userID string, itemOID uint64) error {
	panic("unimplemented")
}

// List implements database.RepeatingItemService
func (*mockRepeatingItemService) List(userID string) ([]dbModels.RepeatingItem, error) {
	panic("unimplemented")
}

// Update implements database.RepeatingItemService
func (*mockRepeatingItemService) Update(userID string, item dbModels.RepeatingItem) error {
	panic("unimplemented")
}

type mockShopService struct{ mock.Mock }

// Create implements database.ShopService
func (*mockShopService) Create(userID string, shop dbModels.Shop) (uint64, error) {
	panic("unimplemented")
}

// Delete implements database.ShopService
func (*mockShopService) Delete(userID string, shopOID uint64) error {
	panic("unimplemented")
}

// Get implements database.ShopService
func (*mockShopService) Get(userID string, shopOID uint64) (*dbModels.Shop, error) {
	panic("unimplemented")
}

// Update implements database.ShopService
func (*mockShopService) Update(userID string, shop dbModels.Shop) error {
	panic("unimplemented")
}

type mockUserService struct{ mock.Mock }

// Create implements database.UserService
func (s *mockUserService) Create(m dbModels.User) error {
	args := s.Called(mock.Anything)
	return args.Error(0)
}

// Get implements database.UserService
func (s *mockUserService) Get(id string) (*dbModels.User, error) {
	args := s.Called(id)
	return args.Get(0).(*dbModels.User), args.Error(1)
}

// GetConfig implements database.UserService
func (s *mockUserService) GetConfig(id string) (dbModels.UserConfig, error) {
	args := s.Called(id)
	return args.Get(0).(dbModels.UserConfig), args.Error(1)
}

// Update implements database.UserService
func (*mockUserService) Update(dbModels.User) error {
	panic("unimplemented")
}

// UpdateConfig implements database.UserService
func (s *mockUserService) UpdateConfig(id string, val dbModels.UserConfig) error {
	args := s.Called(id, val)
	return args.Error(0)
}

type mockDB struct {
	accountService       *mockAccountService
	categoryService      *mockCategoryService
	dailyItemService     *mockDailyItemService
	feeService           *mockFeeService
	repeatingItemService *mockRepeatingItemService
	shopService          *mockShopService
	userService          *mockUserService
}

// Account implements database.DB
func (db *mockDB) Account() database.AccountService {
	return db.accountService
}

// Category implements database.DB
func (db *mockDB) Category() database.CategoryService {
	return db.categoryService
}

// DailyItem implements database.DB
func (db *mockDB) DailyItem() database.DailyItemService {
	return db.dailyItemService
}

// Fee implements database.DB
func (db *mockDB) Fee() database.FeeService {
	return db.feeService
}

// RepeatingItem implements database.DB
func (db *mockDB) RepeatingItem() database.RepeatingItemService {
	return db.repeatingItemService
}

// Shop implements database.DB
func (db *mockDB) Shop() database.ShopService {
	return db.shopService
}

// User implements database.DB
func (db *mockDB) User() database.UserService {
	return db.userService
}

func (db *mockDB) AssertExpectations(t mock.TestingT) bool {
	return db.accountService.AssertExpectations(t) &&
		db.categoryService.AssertExpectations(t) &&
		db.dailyItemService.AssertExpectations(t) &&
		db.feeService.AssertExpectations(t) &&
		db.repeatingItemService.AssertExpectations(t) &&
		db.shopService.AssertExpectations(t) &&
		db.userService.AssertExpectations(t)
}

func NewMockDB() *mockDB {
	return &mockDB{
		accountService:       new(mockAccountService),
		categoryService:      new(mockCategoryService),
		dailyItemService:     new(mockDailyItemService),
		feeService:           new(mockFeeService),
		repeatingItemService: new(mockRepeatingItemService),
		shopService:          new(mockShopService),
		userService:          new(mockUserService),
	}
}
