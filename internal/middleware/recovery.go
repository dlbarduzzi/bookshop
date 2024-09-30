package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/dlbarduzzi/bookshop/internal/logging"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.LoggerFromContext(ctx)

		defer func() {
			if r := recover(); r != nil {
				w.Header().Set("Connection", "close")
				w.WriteHeader(http.StatusInternalServerError)
				log.Error("http handler panic", "panic", r, "stack", string(debug.Stack()))
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
