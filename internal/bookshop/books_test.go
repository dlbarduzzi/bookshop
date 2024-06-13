package bookshop

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
)

func TestCreateBookHandler(t *testing.T) {
	t.Parallel()
	bs := newTestBookshop(t)

	srv := newTestServer(t, bs.Routes())
	defer srv.Close()

	reqBody := `{"title":"Test Book"}`
	code, body := srv.post(t, "/api/v1/books", bytes.NewReader([]byte(reqBody)))

	if code != http.StatusOK {
		t.Errorf("expected status code to be %v; got %v", http.StatusOK, code)
	}

	wantBody := `{Title:Test Book`

	if !strings.Contains(body, wantBody) {
		t.Errorf("expected response body to contain %v; got %v", wantBody, body)
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

	wantBody := `{"book":{"id":1,`

	if !strings.Contains(body, wantBody) {
		t.Errorf("expected response body to contain %v; got %v", wantBody, body)
	}

	code, body = srv.get(t, "/api/v1/books/0")

	if code != http.StatusNotFound {
		t.Errorf("expected status code to be %v; got %v", http.StatusNotFound, code)
	}

	wantBody = "Book with given id was not found."

	if !strings.Contains(body, wantBody) {
		t.Errorf("expected response body to contain %v; got %v", wantBody, body)
	}
}
