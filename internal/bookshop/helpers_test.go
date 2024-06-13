package bookshop

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
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
