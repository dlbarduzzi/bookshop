package bookshop

import (
	"net/http"

	"github.com/dlbarduzzi/bookshop/internal/jsoner"
)

func (bs *Bookshop) healthHandler(w http.ResponseWriter, r *http.Request) {
	data := jsoner.Envelope{
		"status": "healthy",
	}
	if err := jsoner.Marshal(w, data, http.StatusOK, nil); err != nil {
		bs.sendServerError(w, r, err)
		return
	}
}
