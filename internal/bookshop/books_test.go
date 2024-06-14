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

	tests := []struct {
		name     string
		body     string
		wantCode int
		wantBody string
	}{
		{
			name:     "empty_body",
			body:     "",
			wantCode: http.StatusBadRequest,
			wantBody: "request body cannot be empty",
		},
		{
			name:     "empty_body",
			body:     `{"title":""}`,
			wantCode: http.StatusBadRequest,
			wantBody: "At least 1 author is required.",
		},
		{
			name: "valid",
			body: `{"title":"Test Title","authors":["Test Author"],
"published_date":"2024-01-30","page_count":"100 pages","categories":["drama"]}`,
			wantCode: http.StatusCreated,
			wantBody: `"title":"Test Title","authors":`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := newTestServer(t, bs.Routes())
			defer srv.Close()

			code, body := srv.post(t, "/api/v1/books", bytes.NewReader([]byte(tt.body)))

			if code != tt.wantCode {
				t.Errorf("expected status code to be %v; got %v", tt.wantCode, code)
			}

			if !strings.Contains(body, tt.wantBody) {
				t.Errorf("expected response body to contain %v; got %v", tt.wantBody, body)
			}
		})
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
