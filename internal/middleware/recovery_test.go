package middleware

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dlbarduzzi/bookshop/internal/logging"
)

func TestRecovery(t *testing.T) {
	t.Parallel()

	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	ctx := context.Background()
	ctx = logging.LoggerWithContext(ctx, log)

	tests := []struct {
		name    string
		handler http.Handler
		code    int
	}{
		{
			name: "default",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			code: http.StatusOK,
		},
		{
			name: "panic",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("test recovery")
			}),
			code: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			Recovery(tt.handler).ServeHTTP(w, r)

			got := w.Code
			want := tt.code

			if got != want {
				t.Errorf("expected status code to be %v; got %v", want, got)
			}
		})
	}
}
