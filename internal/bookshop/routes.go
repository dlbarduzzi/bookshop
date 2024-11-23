package bookshop

import "net/http"

func (b *Bookshop) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/health", b.healthHandler)
	mux.HandleFunc("GET /api/v1/books", b.listBooksHandler)

	return mux
}
