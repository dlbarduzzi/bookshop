package bookshop

import (
	"bytes"
	"fmt"
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
			name: "force_insert_error",
			body: `{"title":"Force Error","authors":["Test Author"],
"published_date":"2024-01-30","page_count":"100 pages","categories":["drama"]}`,
			wantCode: http.StatusInternalServerError,
			wantBody: "Internal server error.",
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

	tests := []struct {
		name     string
		bookId   int
		wantCode int
		wantBody string
	}{
		{
			name:     "invalid",
			bookId:   0,
			wantCode: http.StatusNotFound,
			wantBody: "Book with given id was not found.",
		},
		{
			name:     "not_found",
			bookId:   1,
			wantCode: http.StatusNotFound,
			wantBody: "Book with given id was not found.",
		},
		{
			name:     "forced_show_error",
			bookId:   2,
			wantCode: http.StatusInternalServerError,
			wantBody: "Internal server error.",
		},
		{
			name:     "valid",
			bookId:   3,
			wantCode: http.StatusOK,
			wantBody: `{"book":{"id":3,"title":"Test Title","authors":`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := newTestServer(t, bs.Routes())
			defer srv.Close()

			code, body := srv.get(t, fmt.Sprintf("/api/v1/books/%d", tt.bookId))

			if code != tt.wantCode {
				t.Errorf("expected status code to be %v; got %v", tt.wantCode, code)
			}

			if !strings.Contains(body, tt.wantBody) {
				t.Errorf("expected response body to contain %v; got %v", tt.wantBody, body)
			}
		})
	}
}
