package guestbook

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListGuestsHandler(t *testing.T) {
	t.Parallel()

	app := newTestGuestbook(t)
	srv := newTestServer(t, app.Routes())
	defer srv.Close()

	code, body := srv.get(t, "/api/v1/guests")
	if code != http.StatusOK {
		t.Errorf("expected status code to be %d; got %d", http.StatusOK, code)
	}

	var res listGuestsResponse

	if err := json.Unmarshal([]byte(body), &res); err != nil {
		t.Fatal(err)
	}

	if len(res.Guests) < 1 {
		t.Errorf("expected list of guests to have at least 1 guest")
	}

	message := res.Guests[0].Message
	wantMessage := "Test Message 1"

	if message != wantMessage {
		t.Errorf("expected guest message to be %s; got %s", wantMessage, message)
	}
}
