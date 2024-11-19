package guestbook

import "net/http"

func (g *Guestbook) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/health", g.healthHandler)

	return mux
}
