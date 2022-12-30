package iris

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/shopspring/decimal"
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

type testScenario[M any] struct {
	Name        string
	RequestBody M
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
