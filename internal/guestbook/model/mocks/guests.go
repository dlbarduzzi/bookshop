package mocks

import (
	"database/sql"

	"github.com/dlbarduzzi/guestbook/internal/guestbook/model"
)

type GuestModel struct {
	DB *sql.DB
}

func (m GuestModel) GetAll() ([]*model.Guest, error) {
	guests := []*model.Guest{
		{
			ID:      1,
			Message: "Test Message 1",
		},
	}
	return guests, nil
}

func (m GuestModel) Insert(guest *model.Guest) error {
	return nil
}
