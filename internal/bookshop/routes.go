package bookshop

import (
	"net/http"

	"github.com/dlbarduzzi/bookshop/internal/middleware"
)

func (bs *Bookshop) Routes() http.Handler {
	mux := http.NewServeMux()

	// Health endpoint.
	mux.HandleFunc("GET /api/v1/health", bs.healthHandler)

	return middleware.Recovery(middleware.RecordRequest(mux))
}
