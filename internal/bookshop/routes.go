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
	mux.HandleFunc("POST /api/v1/books", bs.createBookHandler)
	mux.HandleFunc("GET /api/v1/books/{id}", bs.showBookHandler)
	mux.HandleFunc("PUT /api/v1/books/{id}", bs.updateBookHandler)

	return middleware.Recovery(mux)
}
