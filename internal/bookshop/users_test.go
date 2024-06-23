package bookshop

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
)

func TestRegisterUserHandler(t *testing.T) {
	t.Parallel()
	bs := newTestBookshop(t)

	reqBody := `{"name":"Test User","email":"test.user@email.com","password":"abcd1234"}`

	srv := newTestServer(t, bs.Routes())
	defer srv.Close()

	wantCode := http.StatusCreated

	code, body := srv.post(t, "/api/v1/users", bytes.NewReader([]byte(reqBody)))
	if code != wantCode {
		t.Errorf("expected status code to be %v; got %v", wantCode, code)
	}

	resp := struct {
		User model.User `json:"user"`
	}{}

	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatal(err)
	}

	var wantUserID int64 = 1

	if resp.User.ID != wantUserID {
		t.Errorf("expected user id to be %d; got %d", wantUserID, resp.User.ID)
	}
}
