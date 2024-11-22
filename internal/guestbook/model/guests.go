package model

import (
	"context"
	"database/sql"
	"net"
	"time"
)

type Guest struct {
	ID        int64     `json:"id"`
	IP        net.IP    `json:"ip"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GuestModel struct {
	DB *sql.DB
}

type GuestStore interface {
	GetAll() ([]*Guest, error)
}

func (m GuestModel) GetAll() ([]*Guest, error) {
	query := "SELECT id, ip, message, created_at, updated_at FROM guests LIMIT 10"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	guests := []*Guest{}

	for rows.Next() {
		var guest Guest

		err := rows.Scan(
			&guest.ID,
			&guest.IP,
			&guest.Message,
			&guest.CreatedAt,
			&guest.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		guests = append(guests, &guest)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return guests, nil
}
