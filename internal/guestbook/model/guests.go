package model

import (
	"context"
	"database/sql"
	"time"
)

type Guest struct {
	ID        int64     `json:"id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GuestModel struct {
	DB *sql.DB
}

type GuestStore interface {
	GetAll() ([]*Guest, error)
	Insert(*Guest) error
}

func (m GuestModel) GetAll() ([]*Guest, error) {
	query := "SELECT id, message, created_at, updated_at FROM guests LIMIT 10"

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

func (m GuestModel) Insert(guest *Guest) error {
	query := `
        INSERT INTO guests (message)
        VALUES ($1)
        RETURNING id, created_at, updated_at`

	args := []interface{}{
		guest.Message,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&guest.ID, &guest.CreatedAt, &guest.UpdatedAt,
	)
}
