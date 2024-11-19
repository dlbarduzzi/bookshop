package guestbook

import "net/http"

func (g *Guestbook) healthHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Healthy!"))
}
