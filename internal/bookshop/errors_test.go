package bookshop

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
	bs := newTestBookshop(t)

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	bs.serverError(w, r, fmt.Errorf("test server error"))

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
		Error     string `json:"error"`
		ErrorCode string `json:"error_code"`
	}{}

	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatal(err)
	}

	wantOK := false

	if resp.OK != wantOK {
		t.Errorf("expected ok status to be %v; got %v", wantOK, resp.OK)
	}

	wantError := "Internal server error."

	if resp.Error != wantError {
		t.Errorf("expected error to be %v; got %v", wantError, resp.Error)
	}

	wantErrorCode := ErrCodeServerError

	if resp.ErrorCode != wantErrorCode {
		t.Errorf("expected error code to be %v; got %v", wantErrorCode, resp.ErrorCode)
	}
}

func TestClientError(t *testing.T) {
	t.Parallel()
	bs := newTestBookshop(t)

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	errMessage := "You have made a bad request."

	w := httptest.NewRecorder()
	bs.clientError(w, r, http.StatusBadRequest, errMessage)

	wantCode := http.StatusBadRequest

	if w.Code != wantCode {
		t.Errorf("expected status code to be %v; got %v", wantCode, w.Code)
	}

	body, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	resp := struct {
		OK        bool   `json:"ok"`
		Error     string `json:"error"`
		ErrorCode string `json:"error_code"`
	}{}

	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatal(err)
	}

	wantOK := false

	if resp.OK != wantOK {
		t.Errorf("expected ok status to be %v; got %v", wantOK, resp.OK)
	}

	if resp.Error != errMessage {
		t.Errorf("expected error to be %v; got %v", errMessage, resp.Error)
	}

	if resp.ErrorCode != ErrCodeClientError {
		t.Errorf("expected error code to be %v; got %v", ErrCodeClientError, resp.ErrorCode)
	}
}
