package iris

import (
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LogInAndDoSuiteConfig struct {
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
	return &LogInAndDoSuite{
		db:     NewMockDB(),
		config: cfg,
	}
}

func (s *LogInAndDoSuite) TestDo() {
	resp, err := httpDoWithToken(
		s.config.HTTPRequest.Method,
		s.config.HTTPRequest.URL,
		strings.NewReader(s.config.HTTPRequest.Body),
		s.token,
	)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.config.HTTPExpectation.StatusCode, resp.StatusCode)
	if expected := s.config.HTTPExpectation.BodyJSON; expected != "" {
		assert.JSONEq(s.T(), expected, MustGetResponseBodyJSON(resp.Body))
	}
}

func (s *LogInAndDoSuite) BeforeTest(suiteName, testName string) {
	user := RegisterMockLogIn(s.db)
	if s.config.BeforeTest != nil {
		s.config.BeforeTest(user.ID, s.db)
	}
	myTestServer.db = s.db

	s.token = MustLogIn(user)
}

func (s *LogInAndDoSuite) AfterTest(suiteName, testName string) {
	s.db.AssertExpectations(s.T())
}
