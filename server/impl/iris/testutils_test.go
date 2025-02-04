package iris

import (
	"encoding/base64"
	"io"
	"os"
	"testing"

	"github.com/iris-contrib/httpexpect/v2"
	"github.com/kataras/iris/v12/httptest"
	"github.com/shopspring/decimal"
	"go.uber.org/mock/gomock"

	authV2 "github.com/n101661/maney/pkg/services/auth"
	"github.com/n101661/maney/pkg/utils"
)

type mockService struct {
	auth *authV2.MockService
}

func NewTest(t *testing.T, opts ...utils.Option[options]) (*mockService, *httpexpect.Expect) {
	controller := gomock.NewController(t)

	mockAuth := authV2.NewMockService(controller)

	return &mockService{
		auth: mockAuth,
	}, httptest.New(t, NewServer(&Config{}, mockAuth, opts...).app)
}

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
	myTestServer = NewServer(&Config{
		SecretKey: []byte("maney-secret-key"),
	}, nil)

	go func() {
		if err := myTestServer.ListenAndServe(serverAddr); err != nil {
			os.Exit(1)
		}
	}()

	os.Exit(m.Run())
}

func mustGetResponseBodyJSON(body io.ReadCloser) string {
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

func MustDecimalFromString(s string) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		panic(err)
	}
	return d
}
