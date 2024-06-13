package bookshop

import "net/http"

func (bs *Bookshop) Routes() http.Handler {
	mux := http.NewServeMux()

	// Health endpoint.
	mux.HandleFunc("GET /api/v1/health", bs.healthHandler)

	// Books endpoints.
	mux.HandleFunc("POST /api/v1/books", bs.createBookHandler)
	mux.HandleFunc("GET /api/v1/books/{id}", bs.showBookHandler)

	return mux
}
