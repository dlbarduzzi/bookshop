package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Book struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Authors       []string  `json:"authors"`
	PublishedDate string    `json:"published_date"`
	PageCount     int       `json:"page_count"`
	Categories    []string  `json:"categories"`
	Version       int32     `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type BookModel struct {
	DB *sql.DB
}

type BookStore interface {
	GetAll(Filters) ([]*Book, error)
}

func (m BookModel) GetAll(filters Filters) ([]*Book, error) {
	query := `
        SELECT count(*) OVER(), id, title, authors, TO_CHAR(published_date, 'yyyy-mm-dd'),
        page_count, categories, version, created_at, updated_at
        FROM books LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	args := []interface{}{filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	books := []*Book{}
	count := 0

	for rows.Next() {
		var book Book

		err := rows.Scan(
			&count,
			&book.ID,
			&book.Title,
			pq.Array(&book.Authors),
			&book.PublishedDate,
			&book.PageCount,
			pq.Array(&book.Categories),
			&book.Version,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
