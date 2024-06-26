package bookshop

import (
	"bytes"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model/mocks"
)

func newTestBookshop(t *testing.T) *Bookshop {
	t.Helper()
	return &Bookshop{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		models: mocks.NewModels(&sql.DB{}),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	t.Helper()
	srv := httptest.NewServer(h)
	return &testServer{srv}
}

func (s *testServer) get(t *testing.T, urlPath string) (int, string) {
	res, err := s.Client().Get(s.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)

	return res.StatusCode, string(body)
}

func (s *testServer) post(t *testing.T, urlPath string, jsonBody io.Reader) (int, string) {
	res, err := s.Client().Post(s.URL+urlPath, "", jsonBody)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)

	return res.StatusCode, string(body)
}

func TestPort(t *testing.T) {
	t.Parallel()

	bs := &Bookshop{
		config: &Config{
			Port: 9191,
		},
	}

	wantPort := bs.config.Port

	if bs.Port() != wantPort {
		t.Errorf("expected bookshop port to be %d; got %d", wantPort, bs.Port())
	}
}
