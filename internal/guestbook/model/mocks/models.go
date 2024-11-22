package mocks

import (
	"database/sql"

	"github.com/dlbarduzzi/guestbook/internal/guestbook/model"
)

type Models struct {
	Guests *GuestModel
}

func NewModels(db *sql.DB) model.Models {
	return model.Models{
		Guests: &GuestModel{DB: db},
	}
}
