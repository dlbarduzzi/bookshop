package bookshop

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListBooksHandler(t *testing.T) {
	t.Parallel()

	app := newTestBookshop(t)
	srv := newTestServer(t, app.Routes())
	defer srv.Close()

	code, body := srv.get(t, "/api/v1/books")
	if code != http.StatusOK {
		t.Errorf("expected status code to be %d; got %d", http.StatusOK, code)
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
}
