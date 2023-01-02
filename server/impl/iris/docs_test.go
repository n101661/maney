package iris

// Write test cases like below template for cases that should log in.
//
// func (s *LogInAndDoSuite) TestServer_Xxx() {
// 	s.RunTest(Test{
// 		Name: "your test case name",
// 		BeforeTest: func(userID string, db *mockDB) {
// 			// write down your mock DB behavior.
// 		},
// 		HTTPRequest: HTTPRequest{
// 			Method: "",
// 			URL:    "",
// 			Body:   "", // JSON format.
// 		},
// 		HTTPExpectation: HTTPExpectation{
// 			StatusCode: iris.StatusOK,
// 			BodyJSON:   "", // JSON format if you want assert response body.
// 		},
// 	})
// }
//
// `TestIris_Main` function in `testutils_suite_test.go` is a starting-function
// for those.
