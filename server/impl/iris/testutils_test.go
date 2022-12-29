package iris

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/n101661/maney/database"
	dbModels "github.com/n101661/maney/database/models"
	"github.com/n101661/maney/server/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
)

const serverAddr = "localhost:8080"

var (
	myTestServer   *Server
	myTestPassword = struct {
		Raw       string
		Encrypted []byte
	}{
		Raw:       "my-password",
		Encrypted: mustDecodeBase64("JDJhJDA4JGM1ZlB6WE13MHQzQVBxLkF3QndDYk9uekouR20vSkJOOUg1NUw5LkRBU1R3bkczcVBKWlRD"),
	}
)

func TestMain(m *testing.M) {
	myTestServer = NewServer(Config{
		SecretKey: []byte("maney-secret-key"),
	})

	go func() {
		if err := myTestServer.ListenAndServe(serverAddr); err != nil {
			os.Exit(1)
		}
	}()

	os.Exit(m.Run())
}

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

type testScenario[M any] struct {
	Name        string
	RequestBody M
}

func mustHTTPRequest(method, url string, body interface{}) *http.Request {
	r, err := http.NewRequest(method, url, MustHTTPBody(body))
	if err != nil {
		panic(err)
	}
	return r
}

func httpDoWithToken(method, url string, body interface{}, token string) (*http.Response, error) {
	r := mustHTTPRequest(method, url, body)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	return http.DefaultClient.Do(r)
}

func MustHTTPBody(v interface{}) io.Reader {
	if v == nil {
		return nil
	}

	if r, ok := v.(io.Reader); ok {
		return r
	}

	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(data)
}

func MustGetResponseBody[M any](body io.ReadCloser) *M {
	data, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}

	m := new(M)
	if err := json.Unmarshal(data, m); err != nil {
		panic(err)
	}
	return m
}

func MustGetResponseBodyJSON(body io.ReadCloser) string {
	data, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func mustDecodeBase64(s string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return decoded
}

func MustLogIn(user testUserInfo) (token string) {
	resp, err := http.Post("http://"+serverAddr+"/log-in", "application/json", MustHTTPBody(models.LogInRequestBody{
		ID:       user.ID,
		Password: user.Password.Raw,
	}))
	if err != nil {
		panic(err)
	}

	cookie := resp.Header.Get("Set-Cookie")

	return cookie[6:strings.Index(cookie, ";")]
}

type testUserInfo struct {
	ID       string
	Name     string
	Password struct {
		Raw       string
		Encrypted []byte
	}
}

func RegisterMockLogIn(db *mockDB) testUserInfo {
	user := testUserInfo{
		ID:       "my-id",
		Name:     "tester",
		Password: myTestPassword,
	}

	db.userService.
		On("Get", user.ID).Return(
		&dbModels.User{
			ID:       user.ID,
			Name:     user.Name,
			Password: user.Password.Encrypted,
		}, nil,
	).Once()

	return user
}

func MustDecimalFromString(s string) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		panic(err)
	}
	return d
}
