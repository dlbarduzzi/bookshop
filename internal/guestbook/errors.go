package guestbook

import (
	"log/slog"
	"net/http"

	"github.com/dlbarduzzi/guestbook/internal/jsontil"
)

type serverErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (g *Guestbook) serverError(w http.ResponseWriter, r *http.Request, e error) {
	g.logger.Error(
		e.Error(),
		slog.String("path", r.URL.Path),
		slog.String("method", r.Method),
	)

	res := serverErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error.",
	}

	if err := jsontil.Marshal(w, res, res.Code, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
