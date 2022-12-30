package iris

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	dbModels "github.com/n101661/maney/database/models"
)

type LogInAndDoSuiteConfig struct {
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

	db     *mockDB
	config LogInAndDoSuiteConfig
	token  string
}

func NewLogInAndDoSuite(cfg LogInAndDoSuiteConfig) *LogInAndDoSuite {
	if cfg.Name == "" {
		panic("missing name")
	}

	return &LogInAndDoSuite{
		db:     NewMockDB(),
		config: cfg,
	}
}

func (s *LogInAndDoSuite) TestDo() {
	r, err := http.NewRequest(s.config.HTTPRequest.Method,
		s.config.HTTPRequest.URL,
		strings.NewReader(s.config.HTTPRequest.Body),
	)
	if err != nil {
		panic(err)
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		panic(err)
	}
	assert.Equal(s.T(), s.config.HTTPExpectation.StatusCode, resp.StatusCode, s.config.Name)
	if expected := s.config.HTTPExpectation.BodyJSON; expected != "" {
		assert.JSONEq(s.T(), expected, mustGetResponseBodyJSON(resp.Body), s.config.Name)
	}
}

func (s *LogInAndDoSuite) BeforeTest(suiteName, testName string) {
	var (
		userID       = "my-id"
		userName     = "tester"
		userPassword = myTestPassword
	)

	s.db.userService.On("Get", userID).Return(
		&dbModels.User{
			ID:       userID,
			Name:     userName,
			Password: userPassword.Encrypted,
		}, nil,
	).Once()
	if s.config.BeforeTest != nil {
		s.config.BeforeTest(userID, s.db)
	}
	myTestServer.db = s.db

	{ // log in
		resp, err := http.Post("http://"+serverAddr+"/log-in", "application/json", strings.NewReader(fmt.Sprintf(`{
			"id": "%s",
			"password": "%s"
		}`, userID, userPassword.Raw)))
		if err != nil {
			panic(err)
		}

		cookie := resp.Header.Get("Set-Cookie")

		s.token = cookie[6:strings.Index(cookie, ";")]
	}
}

func (s *LogInAndDoSuite) AfterTest(suiteName, testName string) {
	s.db.AssertExpectations(s.T())
}
