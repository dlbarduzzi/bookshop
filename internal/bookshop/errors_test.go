package bookshop

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendServerError(t *testing.T) {
	t.Parallel()
	bs := newTestBookshop(t)

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	bs.sendServerError(w, r, fmt.Errorf("test server error"))

	wantCode := http.StatusInternalServerError

	if w.Code != wantCode {
		t.Errorf("expected status code to be %v; got %v", wantCode, w.Code)
	}

	body, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	resp := struct {
		OK        bool   `json:"ok"`
		Message   string `json:"message"`
		ErrorCode string `json:"errorCode"`
	}{}

	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatal(err)
	}

	wantOK := false

	if resp.OK != wantOK {
		t.Errorf("expected ok status to be %v; got %v", wantOK, resp.OK)
	}

	wantError := "Something went wrong while processing your request."

	if resp.Message != wantError {
		t.Errorf("expected error to be %v; got %v", wantError, resp.Message)
	}

	wantErrorCode := ErrCodeServerError

	if resp.ErrorCode != wantErrorCode {
		t.Errorf("expected error code to be %v; got %v", wantErrorCode, resp.ErrorCode)
	}
}
