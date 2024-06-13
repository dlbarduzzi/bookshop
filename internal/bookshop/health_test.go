package bookshop

import (
	"net/http"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	t.Parallel()
	bs := newTestBookshop(t)

	srv := newTestServer(t, bs.Routes())
	defer srv.Close()

	code, body := srv.get(t, "/api/v1/health")

	if code != http.StatusOK {
		t.Errorf("expected status code to be %v; got %v", http.StatusOK, code)
	}

	wantBody := `{"status":"ok"}`

	if body != wantBody {
		t.Errorf("expected response body to be %v; got %v", wantBody, body)
	}
}
