package guestbook

import (
	"net/http"

	"github.com/dlbarduzzi/guestbook/internal/guestbook/model"
	"github.com/dlbarduzzi/guestbook/internal/jsontil"
)

type listGuestsResponse struct {
	Code   int            `json:"code"`
	Guests []*model.Guest `json:"guests"`
}

func (g *Guestbook) listGuestsHandler(w http.ResponseWriter, r *http.Request) {
	guests, err := g.models.Guests.GetAll()
	if err != nil {
		g.serverError(w, r, err)
		return
	}

	res := listGuestsResponse{
		Code:   http.StatusOK,
		Guests: guests,
	}

	if err := jsontil.Marshal(w, res, res.Code, nil); err != nil {
		g.serverError(w, r, err)
		return
	}
}

type createGuestResponse struct {
	Code  int          `json:"code"`
	Guest *model.Guest `json:"guest"`
}

func (g *Guestbook) createGuestHandler(w http.ResponseWriter, r *http.Request) {
	guest := &model.Guest{
		Message: "A sample message",
	}

	if err := g.models.Guests.Insert(guest); err != nil {
		g.serverError(w, r, err)
		return
	}

	res := createGuestResponse{
		Code:  http.StatusOK,
		Guest: guest,
	}

	if err := jsontil.Marshal(w, res, res.Code, nil); err != nil {
		g.serverError(w, r, err)
		return
	}
}
