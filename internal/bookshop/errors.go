package bookshop

import (
	"log/slog"
	"net/http"

	"github.com/dlbarduzzi/bookshop/internal/jsoner"
)

var ErrCodeServerError = "server-error"

func (bs *Bookshop) serverError(w http.ResponseWriter, r *http.Request, err error) {
	bs.logger.Error(err.Error(),
		slog.String("method", r.Method),
		slog.String("url", r.URL.RequestURI()),
	)
	data := jsoner.Envelope{
		"ok":         false,
		"error":      "Internal server error.",
		"error_code": ErrCodeServerError,
	}
	if err := jsoner.Marshal(w, data, http.StatusInternalServerError, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
