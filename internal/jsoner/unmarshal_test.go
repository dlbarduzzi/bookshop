package jsoner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	code, err := unmarshalTestHelper(t, `{"foo":"bar"}`)
	if err != nil {
		t.Errorf("expected error to be nil; got %v", err)
	}

	if code != http.StatusOK {
		t.Errorf("expected status code to be %d; got %d", http.StatusOK, code)
	}
}

func TestInvalidRequestBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		wantErr  string
		wantCode int
	}{
		{
			name:     "malformed_at_position",
			input:    `{"foo":"bar",}`,
			wantErr:  "request body contains malformed json at position 14",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "malformed_json",
			input:    `{"foo":"bar"`,
			wantErr:  "request body contains malformed json",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "incorrect_type_at_position",
			input:    `{"foo":1001}`,
			wantErr:  "request body contains incorrect json type \"foo\" at position 11",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "empty_field",
			input:    `{"":1001}`,
			wantErr:  "request body contains empty json field name",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "unknown_field",
			input:    `{"baz":"bar"}`,
			wantErr:  "request body contains unknown json field \"baz\"",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "empty_body",
			input:    ``,
			wantErr:  "request body cannot be empty",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "multiple_bodies",
			input:    `{"foo":"bar"}{"foo":"bar"}`,
			wantErr:  "request body cannot have more than 1 json object",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, err := unmarshalTestHelper(t, tt.input)

			if err == nil || err.Error() != tt.wantErr {
				t.Errorf("expected error to be %v; got %v", tt.wantErr, err)
			}

			if code != tt.wantCode {
				t.Errorf("expected status code to be %d; got %d", tt.wantCode, code)
			}
		})
	}
}

func TestLargeRequestBody(t *testing.T) {
	t.Parallel()

	data := make(map[string]string, 1)
	data["foo"] = strings.Repeat("0", maxBodyBytes+10)

	input, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	code, err := unmarshalTestHelper(t, string(input))

	wantErr := fmt.Sprintf("request body cannot be larger than %d bytes", maxBodyBytes)
	wantCode := http.StatusRequestEntityTooLarge

	if err == nil || err.Error() != wantErr {
		t.Errorf("expected error to be %v; got %v", wantErr, err)
	}

	if code != wantCode {
		t.Errorf("expected status code to be %d; got %d", wantCode, code)
	}
}

type unmarshalTestFoo struct {
	Foo string `json:"foo"`
}

func unmarshalTestHelper(t *testing.T, input string) (int, error) {
	t.Helper()

	body := io.NopCloser(bytes.NewReader([]byte(input)))

	r := httptest.NewRequest("POST", "/", body)
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	data := &unmarshalTestFoo{}

	return Unmarshal(w, r, data)
}

func TestServerError(t *testing.T) {
	code, err := serverError(errors.New("forced error"))

	wantErr := "failed to decode json request body; forced error"
	wantCode := http.StatusInternalServerError

	if err == nil || err.Error() != wantErr {
		t.Errorf("expected error to be %v; got %v", wantErr, err)
	}

	if code != wantCode {
		t.Errorf("expected status code to be %d; got %d", wantCode, code)
	}
}
