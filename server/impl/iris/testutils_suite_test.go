package iris

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	dbModels "github.com/n101661/maney/database/models"
)

type Test struct {
	Name string

	BeforeTest func(userID string, db *mockDB)

	HTTPRequest     HTTPRequest
	HTTPExpectation HTTPExpectation
}

type HTTPRequest struct {
	Method string
	URL    string
	Body   string
}

type HTTPExpectation struct {
	StatusCode int
	BodyJSON   string
}

type LogInAndDoSuite struct {
	suite.Suite

	userID string
	token  string
}

func NewLogInAndDoSuite(t *testing.T) *LogInAndDoSuite {
	var (
		userID       = "my-id"
		userName     = "tester"
		userPassword = myTestPassword
	)

	db := NewMockDB()
	db.userService.On("Get", userID).Return(
		&dbModels.User{
			ID:       userID,
			Name:     userName,
			Password: userPassword.Encrypted,
		}, nil,
	).Once()
	myTestServer.db = db

	{ // log in
		resp, err := http.Post("http://"+serverAddr+"/log-in", "application/json", strings.NewReader(fmt.Sprintf(`{
			"id": "%s",
			"password": "%s"
		}`, userID, userPassword.Raw)))
		if err != nil {
			panic(err)
		}

		cookie := resp.Header.Get("Set-Cookie")

		db.AssertExpectations(t)

		return &LogInAndDoSuite{
			userID: userID,
			token:  cookie[6:strings.Index(cookie, ";")],
		}
	}
}

func (s *LogInAndDoSuite) RunTest(test Test) {
	if test.Name == "" {
		panic("missing name")
	}

	db := NewMockDB()
	if test.BeforeTest != nil {
		test.BeforeTest(s.userID, db)
	}
	myTestServer.db = db

	r, err := http.NewRequest(test.HTTPRequest.Method,
		test.HTTPRequest.URL,
		strings.NewReader(test.HTTPRequest.Body),
	)
	if err != nil {
		panic(err)
	}

	switch strings.ToLower(test.HTTPRequest.Method) {
	case "post", "put":
		r.Header.Set("Content-Type", "application/json")
	}
	r.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		panic(err)
	}
	assert.Equal(s.T(), test.HTTPExpectation.StatusCode, resp.StatusCode, test.Name)
	if expected := test.HTTPExpectation.BodyJSON; expected != "" {
		assert.JSONEq(s.T(), expected, mustGetResponseBodyJSON(resp.Body), test.Name)
	}

	db.AssertExpectations(s.T())
}

func TestIris_Main(t *testing.T) {
	suite.Run(t, NewLogInAndDoSuite(t))
}
