package guestbook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerError(t *testing.T) {
	t.Parallel()

	app := newTestGuestbook(t)

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	app.serverError(w, r, fmt.Errorf("test server error"))

	body, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	var res serverErrorResponse

	if err := json.Unmarshal(body, &res); err != nil {
		t.Fatal(err)
	}

	wantCode := http.StatusInternalServerError
	if res.Code != wantCode {
		t.Errorf("expected status code to be %d; got %d", wantCode, res.Code)
	}

	wantMessage := "Internal server error."
	if res.Message != wantMessage {
		t.Errorf("expected error message to be %s; got %s", wantMessage, res.Message)
	}
}
