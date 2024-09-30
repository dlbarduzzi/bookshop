package bookshop

import (
	"log/slog"
	"net/http"

	"github.com/dlbarduzzi/bookshop/internal/jsoner"
)

const ErrCodeServerError = "internal-server-error"

func (bs *Bookshop) sendServerError(w http.ResponseWriter, r *http.Request, err error) {
	bs.logger.Error(err.Error(),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)
	data := jsoner.Envelope{
		"ok":        false,
		"message":   "Something went wrong while processing your request.",
		"errorCode": ErrCodeServerError,
	}
	if err := jsoner.Marshal(w, data, http.StatusInternalServerError, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
