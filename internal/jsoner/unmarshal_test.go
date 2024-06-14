package jsoner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	err := unmarshalTestHelper(t, `{"foo":"bar"}`)
	if err != nil {
		t.Errorf("expected error to be nil; got %v", err)
	}
}

func TestInvalidRequestBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name:    "malformed_at_position",
			input:   `{"foo":"bar",}`,
			wantErr: "request body contains malformed json at position 14",
		},
		{
			name:    "malformed_json",
			input:   `{"foo":"bar"`,
			wantErr: "request body contains malformed json",
		},
		{
			name:    "incorrect_type_at_position",
			input:   `{"foo":1001}`,
			wantErr: "request body contains incorrect json type \"foo\" at position 11",
		},
		{
			name:    "empty_field",
			input:   `{"":1001}`,
			wantErr: "request body contains empty json field name",
		},
		{
			name:    "unknown_field",
			input:   `{"baz":"bar"}`,
			wantErr: "request body contains unknown json field \"baz\"",
		},
		{
			name:    "empty_body",
			input:   ``,
			wantErr: "request body cannot be empty",
		},
		{
			name:    "multiple_bodies",
			input:   `{"foo":"bar"}{"foo":"bar"}`,
			wantErr: "request body cannot have more than 1 json object",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := unmarshalTestHelper(t, tt.input)
			if err == nil || err.Error() != tt.wantErr {
				t.Errorf("expected error to be %v; got %v", tt.wantErr, err)
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

	err = unmarshalTestHelper(t, string(input))
	wantErr := fmt.Sprintf("request body cannot be larger than %d bytes", maxBodyBytes)

	if err == nil || err.Error() != wantErr {
		t.Errorf("expected error to be %v; got %v", wantErr, err)
	}
}

type unmarshalTestFoo struct {
	Foo string `json:"foo"`
}

func unmarshalTestHelper(t *testing.T, input string) error {
	t.Helper()

	body := io.NopCloser(bytes.NewReader([]byte(input)))

	r := httptest.NewRequest("POST", "/", body)
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	data := &unmarshalTestFoo{}

	return Unmarshal(w, r, data)
}
