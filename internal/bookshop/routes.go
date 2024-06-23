package bookshop

import (
	"net/http"

	"github.com/dlbarduzzi/bookshop/internal/middleware"
)

func (bs *Bookshop) Routes() http.Handler {
	mux := http.NewServeMux()

	// Health endpoint.
	mux.HandleFunc("GET /api/v1/health", bs.healthHandler)

	// Books endpoints.
	mux.HandleFunc("GET /api/v1/books", bs.listBookHandler)
	mux.HandleFunc("POST /api/v1/books", bs.createBookHandler)
	mux.HandleFunc("GET /api/v1/books/{id}", bs.showBookHandler)
	mux.HandleFunc("PATCH /api/v1/books/{id}", bs.updateBookHandler)
	mux.HandleFunc("DELETE /api/v1/books/{id}", bs.deleteBookHandler)

	// Users endpoints.
	mux.HandleFunc("POST /api/v1/users", bs.registerUserHandler)
	mux.HandleFunc("GET /api/v1/users/verify-email", bs.verifyEmailUserHandler)

	// Should read these values from env variables.
	limiter := middleware.Limiter{
		RPS:     2,
		Burst:   4,
		Enabled: true,
	}

	return middleware.Recovery(limiter.RateLimit(mux))
}
