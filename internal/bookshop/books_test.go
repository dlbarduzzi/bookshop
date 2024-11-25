package bookshop

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListBooksHandler(t *testing.T) {
	t.Parallel()
	app := newTestBookshop(t)

	testCases := []struct {
		name     string
		query    string
		wantCode int
	}{
		{
			name:     "invalid page parameter",
			query:    "?page=not_an_integer",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "valid page parameter",
			query:    "?page=10",
			wantCode: http.StatusOK,
		},
		{
			name:     "success",
			query:    "",
			wantCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := newTestServer(t, app.Routes())
			defer srv.Close()

			code, body := srv.get(t, "/api/v1/books"+tc.query)
			if code != tc.wantCode {
				t.Errorf("expected status code to be %d; got %d", tc.wantCode, code)
			}

			if tc.wantCode != http.StatusOK {
				var res validationErrorResponse

				if err := json.Unmarshal([]byte(body), &res); err != nil {
					t.Fatal(err)
				}

				wantMessage := "Input validation error."
				if res.Message != wantMessage {
					t.Errorf("expected error message to be %s; got %s", wantMessage, res.Message)
				}
				return
			}

			var res listBooksResponse

			if err := json.Unmarshal([]byte(body), &res); err != nil {
				t.Fatal(err)
			}

			if len(res.Books) < 1 {
				t.Errorf("expected list of books to have at least 1 book")
			}

			bookTitle := res.Books[0].Title
			wantBookTitle := "Test Book 1"

			if bookTitle != wantBookTitle {
				t.Errorf("expected book title to be %s; got %s", wantBookTitle, bookTitle)
			}
		})
	}
}
