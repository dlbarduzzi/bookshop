package bookshop

import (
	"net/http"
	"testing"
)

func TestCreateBookHandler(t *testing.T) {
	t.Parallel()
	bs := newTestBookshop(t)

	srv := newTestServer(t, bs.Routes())
	defer srv.Close()

	code, body := srv.post(t, "/api/v1/books", nil)

	if code != http.StatusOK {
		t.Errorf("expected status code to be %v; got %v", http.StatusOK, code)
	}

	wantBody := "create a new book"

	if body != wantBody {
		t.Errorf("expected response body to be %v; got %v", wantBody, body)
	}
}

func TestShowBookHandler(t *testing.T) {
	t.Parallel()
	bs := newTestBookshop(t)

	srv := newTestServer(t, bs.Routes())
	defer srv.Close()

	code, body := srv.get(t, "/api/v1/books/1")

	if code != http.StatusOK {
		t.Errorf("expected status code to be %v; got %v", http.StatusOK, code)
	}

	wantBody := "show book with id 1"

	if body != wantBody {
		t.Errorf("expected response body to be %v; got %v", wantBody, body)
	}

	code, body = srv.get(t, "/api/v1/books/0")

	if code != http.StatusNotFound {
		t.Errorf("expected status code to be %v; got %v", http.StatusNotFound, code)
	}

	wantBody = "404 page not found"

	if body != wantBody {
		t.Errorf("expected response body to be %v; got %v", wantBody, body)
	}
}
