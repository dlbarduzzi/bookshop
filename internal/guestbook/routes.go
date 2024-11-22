package guestbook

import "net/http"

func (g *Guestbook) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/health", g.healthHandler)

	mux.HandleFunc("GET /api/v1/guests", g.listGuestsHandler)
	mux.HandleFunc("POST /api/v1/guests", g.createGuestHandler)

	return mux
}
