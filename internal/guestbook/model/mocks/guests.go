package mocks

import (
	"database/sql"
	"net"

	"github.com/dlbarduzzi/guestbook/internal/guestbook/model"
)

type GuestModel struct {
	DB *sql.DB
}

func (m GuestModel) GetAll() ([]*model.Guest, error) {
	guests := []*model.Guest{
		{
			ID:      1,
			IP:      net.IPv4(10, 0, 0, 1),
			Message: "Test Message 1",
		},
	}
	return guests, nil
}
