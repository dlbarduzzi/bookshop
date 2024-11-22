package guestbook

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	t.Parallel()

	app := newTestGuestbook(t)
	srv := newTestServer(t, app.Routes())
	defer srv.Close()

	code, body := srv.get(t, "/api/v1/health")
	if code != http.StatusOK {
		t.Errorf("expected status code to be %d; got %d", http.StatusOK, code)
	}

	var res healthResponse

	if err := json.Unmarshal([]byte(body), &res); err != nil {
		t.Fatal(err)
	}

	wantMessage := "API is healthy."
	if res.Message != wantMessage {
		t.Errorf("expected response message to be %s; got %s", wantMessage, res.Message)
	}
}
