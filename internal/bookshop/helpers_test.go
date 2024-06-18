package bookshop

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/dlbarduzzi/bookshop/internal/validator"
)

func TestReadIDParam(t *testing.T) {
	t.Parallel()

	bs := newTestBookshop(t)

	tests := []struct {
		name    string
		value   string
		wantErr string
	}{
		{
			name:    "valid",
			value:   "123",
			wantErr: "",
		},
		{
			name:    "less_than_one",
			value:   "0",
			wantErr: "invalid id parameter",
		},
		{
			name:    "not_an_integer",
			value:   "abc",
			wantErr: "invalid id parameter",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.SetPathValue("id", tt.value)

			id, err := bs.readIDParam(r)

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("expected error to be nil; got %v", err)
				}

				if fmt.Sprintf("%d", id) != tt.value {
					t.Errorf("expected id param to be %d; got %d", 1, id)
				}
			}

			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("expected error to be %v; got %v", tt.wantErr, err)
				}

				if id != 0 {
					t.Errorf("expected id param to be 0; got %d", id)
				}
			}
		})
	}
}

func TestReadString(t *testing.T) {
	t.Parallel()

	bs := newTestBookshop(t)

	tests := []struct {
		name         string
		key          string
		value        string
		wantValue    string
		defaultValue string
	}{
		{
			name:         "default",
			key:          "title",
			value:        "",
			wantValue:    "test",
			defaultValue: "test",
		},
		{
			name:         "query_value",
			key:          "title",
			value:        "value",
			wantValue:    "value",
			defaultValue: "test",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(
				http.MethodGet, fmt.Sprintf("/test?%s=%s", tt.key, tt.value), nil)

			val := bs.readString(r.URL.Query(), tt.key, tt.defaultValue)
			if val != tt.wantValue {
				t.Errorf("expected key value to be %s; got %s", tt.wantValue, val)
			}
		})
	}
}

func TestReadCSV(t *testing.T) {
	t.Parallel()

	bs := newTestBookshop(t)

	tests := []struct {
		name         string
		key          string
		value        string
		wantValue    []string
		defaultValue []string
	}{
		{
			name:         "default",
			key:          "title",
			value:        "",
			wantValue:    []string{"test"},
			defaultValue: []string{"test"},
		},
		{
			name:         "query_value",
			key:          "title",
			value:        "value",
			wantValue:    []string{"value"},
			defaultValue: []string{"test"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(
				http.MethodGet, fmt.Sprintf("/test?%s=%s", tt.key, tt.value), nil)

			val := bs.readCSV(r.URL.Query(), tt.key, tt.defaultValue)
			if !reflect.DeepEqual(val, tt.wantValue) {
				t.Errorf("expected key value to be %s; got %s", tt.wantValue, val)
			}
		})
	}
}

func TestReadInt(t *testing.T) {
	t.Parallel()

	bs := newTestBookshop(t)

	tests := []struct {
		name         string
		key          string
		value        string
		wantValue    int
		defaultValue int
	}{
		{
			name:         "default",
			key:          "year",
			value:        "",
			wantValue:    2024,
			defaultValue: 2024,
		},
		{
			name:         "default_not_an_int",
			key:          "year",
			value:        "abc",
			wantValue:    2024,
			defaultValue: 2024,
		},
		{
			name:         "query_value",
			key:          "year",
			value:        "1994",
			wantValue:    1994,
			defaultValue: 2024,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(
				http.MethodGet, fmt.Sprintf("/test?%s=%s", tt.key, tt.value), nil)

			v := validator.NewValidator()

			val := bs.readInt(r.URL.Query(), tt.key, tt.defaultValue, v)
			if val != tt.wantValue {
				t.Errorf("expected key value to be %d; got %d", tt.wantValue, val)
			}
		})
	}
}
