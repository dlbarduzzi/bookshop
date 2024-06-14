package bookshop

import (
	"log/slog"
	"net/http"

	"github.com/dlbarduzzi/bookshop/internal/jsoner"
)

var (
	ErrCodeServerError     = "server-error"
	ErrCodeClientError     = "client-error"
	ErrCodeValidationError = "validation-error"
)

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

func (bs *Bookshop) clientError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	data := jsoner.Envelope{
		"ok":         false,
		"error":      msg,
		"error_code": ErrCodeClientError,
	}
	if err := jsoner.Marshal(w, data, code, nil); err != nil {
		bs.serverError(w, r, err)
		return
	}
}

func (bs *Bookshop) validationError(w http.ResponseWriter, r *http.Request, errs any) {
	data := jsoner.Envelope{
		"ok":         false,
		"error":      errs,
		"error_code": ErrCodeValidationError,
	}
	if err := jsoner.Marshal(w, data, http.StatusBadRequest, nil); err != nil {
		bs.serverError(w, r, err)
		return
	}
}
