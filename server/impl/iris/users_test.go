package iris

import (
	"errors"

	"github.com/kataras/iris/v12"

	dbModels "github.com/n101661/maney/database/models"
)

func (s *LogInAndDoSuite) TestServer_UpdateConfig() {
	const addr = "http://" + serverAddr + "/users/config"

	{
		s.RunTest(Test{
			Name: "failed to update config(unexpected error)",
			BeforeTest: func(userID string, db *mockDB) {
				db.userService.On("UpdateConfig", userID, dbModels.UserConfig{
					CompareItemsInDifferentShop: false,
					CompareItemsInSameShop:      true,
				}).Return(errors.New("unexpected error")).Once()
			},
			HTTPRequest: HTTPRequest{
				Method: "PUT",
				URL:    addr,
				Body: `{
					"compare_items_in_different_shop": false,
					"compare_items_in_same_shop": true
				}`,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusInternalServerError,
			},
		})
	}
	{
		s.RunTest(Test{
			Name: "update config successful",
			BeforeTest: func(userID string, db *mockDB) {
				db.userService.On("UpdateConfig", userID, dbModels.UserConfig{
					CompareItemsInDifferentShop: false,
					CompareItemsInSameShop:      true,
				}).Return(nil).Once()
			},
			HTTPRequest: HTTPRequest{
				Method: "PUT",
				URL:    addr,
				Body: `{
					"compare_items_in_different_shop": false,
					"compare_items_in_same_shop": true
				}`,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusOK,
			},
		})
	}
}

func (s *LogInAndDoSuite) TestServer_GetConfig() {
	const addr = "http://" + serverAddr + "/users/config"

	{
		s.RunTest(Test{
			Name: "failed to get config(unexpected error)",
			BeforeTest: func(userID string, db *mockDB) {
				db.userService.On("GetConfig", userID).Return(dbModels.UserConfig{}, errors.New("unexpected error")).Once()
			},
			HTTPRequest: HTTPRequest{
				Method: "GET",
				URL:    addr,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusInternalServerError,
			},
		})
	}
	{
		s.RunTest(Test{
			Name: "get config successful",
			BeforeTest: func(userID string, db *mockDB) {
				db.userService.On("GetConfig", userID).Return(dbModels.UserConfig{
					CompareItemsInDifferentShop: true,
					CompareItemsInSameShop:      false,
				}, nil).Once()
			},
			HTTPRequest: HTTPRequest{
				Method: "GET",
				URL:    addr,
			},
			HTTPExpectation: HTTPExpectation{
				StatusCode: iris.StatusOK,
				BodyJSON: `{
					"compare_items_in_different_shop": true,
					"compare_items_in_same_shop": false
				}`,
			},
		})
	}
}
