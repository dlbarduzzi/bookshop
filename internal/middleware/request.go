package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/dlbarduzzi/bookshop/internal/logging"
)

type responseWriter struct {
	statusCode int
	http.ResponseWriter
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func RecordRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ctx := r.Context()
		log := logging.LoggerFromContext(ctx)

		rw := &responseWriter{
			statusCode:     http.StatusOK,
			ResponseWriter: w,
		}

		next.ServeHTTP(rw, r)

		log.Info("request details",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("statusCode", rw.statusCode),
			slog.Duration("duration", time.Since(start)))
	})
}
