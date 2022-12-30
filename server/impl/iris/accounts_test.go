package iris

import (
	"errors"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	"github.com/n101661/maney/database"
	dbModels "github.com/n101661/maney/database/models"
)

func TestServer_CreateAccount(t *testing.T) {
	const addr = "http://" + serverAddr + "/users/accounts"

	{
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			Name: "missing name",
			HTTPRequest: HTTPRequest{
				Method: "POST",
				URL:    addr,
				Body: `{
					"name": "",
					"icon_oid": "0",
					"initial_balance": "0"
				}`,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusBadRequest,
			},
		}))
	}
	{
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			Name: "the account has already existed",
			BeforeTest: func(userID string, db *mockDB) {
				db.accountService.On("Create", userID, dbModels.AssetAccount{
					Name:           "account-name",
					IconOID:        0,
					InitialBalance: decimal.Zero,
					Balance:        decimal.Zero,
				}).Return(uint64(0), database.ErrResourceExisted).Once()
			},
			HTTPRequest: HTTPRequest{
				Method: "POST",
				URL:    addr,
				Body: `{
					"name": "account-name",
					"icon_oid": "0",
					"initial_balance": "0"
				}`,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusConflict,
			},
		}))
	}
	{
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			Name: "failed to create account(unexpected error)",
			BeforeTest: func(userID string, db *mockDB) {
				db.accountService.On("Create", userID, dbModels.AssetAccount{
					Name:           "account-name",
					IconOID:        0,
					InitialBalance: decimal.Zero,
					Balance:        decimal.Zero,
				}).Return(uint64(0), errors.New("unexpected error")).Once()
			},
			HTTPRequest: HTTPRequest{
				Method: "POST",
				URL:    addr,
				Body: `{
					"name": "account-name",
					"icon_oid": "0",
					"initial_balance": "0"
				}`,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusInternalServerError,
			},
		}))
	}
	{
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			Name: "create account successful",
			BeforeTest: func(userID string, db *mockDB) {
				db.accountService.On("Create", userID, dbModels.AssetAccount{
					Name:           "account-name",
					IconOID:        0,
					InitialBalance: decimal.Zero,
					Balance:        decimal.Zero,
				}).Return(uint64(0), nil).Once()
			},
			HTTPRequest: HTTPRequest{
				Method: "POST",
				URL:    addr,
				Body: `{
					"name": "account-name",
					"icon_oid": "0",
					"initial_balance": "0"
				}`,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusOK,
			},
		}))
	}
}

func TestServer_ListAccounts(t *testing.T) {
	const addr = "http://" + serverAddr + "/users/accounts"

	{
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			Name: "failed to get accounts(unexpected error)",
			BeforeTest: func(userID string, db *mockDB) {
				db.accountService.On("List", userID).Return(([]dbModels.AssetAccount)(nil), errors.New("unexpected error")).Once()
			},
			HTTPRequest: HTTPRequest{
				Method: "GET",
				URL:    addr,
				Body: `{
					"name": "account-name",
					"icon_oid": "0",
					"initial_balance": "0"
				}`,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusInternalServerError,
			},
		}))
	}
	{
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			Name: "get accounts successful",
			BeforeTest: func(userID string, db *mockDB) {
				db.accountService.On("List", userID).Return([]dbModels.AssetAccount{{
					OID:            0,
					Name:           "account-1",
					IconOID:        0,
					InitialBalance: decimal.Zero,
					Balance:        decimal.NewFromFloat32(12.3),
				}, {
					OID:            1,
					Name:           "account-2",
					IconOID:        2,
					InitialBalance: decimal.NewFromFloat32(12.3),
					Balance:        decimal.NewFromFloat32(4.56),
				}}, nil).Once()
			},
			HTTPRequest: HTTPRequest{
				Method: "GET",
				URL:    addr,
				Body: `{
					"name": "account-name",
					"icon_oid": "0",
					"initial_balance": "0"
				}`,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusOK,
				BodyJSON: `[
					{
						"oid": "0",
						"name": "account-1",
						"icon_oid": "0",
						"initial_balance": "0",
						"balance": "12.3"
					}, {
						"oid": "1",
						"name": "account-2",
						"icon_oid": "2",
						"initial_balance": "12.3",
						"balance": "4.56"
					}
				]`,
			},
		}))
	}
}

func TestServer_UpdateAccount(t *testing.T) {
	const addr = "http://" + serverAddr + "/users/accounts"

	{
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			Name: "missing account oid",
			HTTPRequest: HTTPRequest{
				Method: "PUT",
				URL:    addr + "/",
				Body: `{
					"name": "my-account",
					"icon_oid": "10",
					"initial_balance": "7.89"
				}`,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusBadRequest,
			},
		}))
	}
	{
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			Name: "invalid account oid",
			HTTPRequest: HTTPRequest{
				Method: "PUT",
				URL:    addr + "/abc",
				Body: `{
					"name": "my-account",
					"icon_oid": "10",
					"initial_balance": "7.89"
				}`,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusBadRequest,
			},
		}))
	}
	{
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			Name: "failed to update account(unexpected error)",
			BeforeTest: func(userID string, db *mockDB) {
				db.accountService.On("Update", userID, dbModels.AssetAccount{
					OID:            0,
					Name:           "my-account",
					IconOID:        10,
					InitialBalance: MustDecimalFromString("7.89"),
				}).Return(errors.New("unexpected error"))
			},
			HTTPRequest: HTTPRequest{
				Method: "PUT",
				URL:    addr + "/0",
				Body: `{
					"name": "my-account",
					"icon_oid": "10",
					"initial_balance": "7.89"
				}`,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusInternalServerError,
			},
		}))
	}
	{
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			Name: "update account successful",
			BeforeTest: func(userID string, db *mockDB) {
				db.accountService.On("Update", userID, dbModels.AssetAccount{
					OID:            0,
					Name:           "my-account",
					IconOID:        10,
					InitialBalance: MustDecimalFromString("7.89"),
				}).Return(nil)
			},
			HTTPRequest: HTTPRequest{
				Method: "PUT",
				URL:    addr + "/0",
				Body: `{
					"name": "my-account",
					"icon_oid": "10",
					"initial_balance": "7.89"
				}`,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusOK,
			},
		}))
	}
}

func TestServer_DeleteAccount(t *testing.T) {
	const addr = "http://" + serverAddr + "/users/accounts"

	{
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			Name: "missing account oid",
			HTTPRequest: HTTPRequest{
				Method: "DELETE",
				URL:    addr + "/",
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusBadRequest,
			},
		}))
	}
	{
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			Name: "invalid account oid",
			HTTPRequest: HTTPRequest{
				Method: "DELETE",
				URL:    addr + "/abc",
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusBadRequest,
			},
		}))
	}
	{
		suite.Run(t, NewLogInAndDoSuite(LogInAndDoSuiteConfig{
			Name: "delete account successful",
			BeforeTest: func(userID string, db *mockDB) {
				db.accountService.On("Delete", userID, uint64(99)).Return(nil)
			},
			HTTPRequest: HTTPRequest{
				Method: "DELETE",
				URL:    addr + "/99",
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusOK,
			},
		}))
	}
}
