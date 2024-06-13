package bookshop

import (
	"net/http"

	"github.com/dlbarduzzi/bookshop/internal/jsoner"
)

func (bs *Bookshop) healthHandler(w http.ResponseWriter, r *http.Request) {
	data := jsoner.Envelope{
		"status":  "ok",
		"version": version,
	}
	if err := jsoner.Marshal(w, data, http.StatusOK, nil); err != nil {
		bs.serverError(w, r, err)
		return
	}
}
