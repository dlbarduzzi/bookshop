package model

import (
	"context"
	"database/sql"
	"fmt"
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
	GetAll(string, []string, Filters) ([]*Book, Metadata, error)
}

func (m BookModel) GetAll(
	title string,
	categories []string,
	filters Filters,
) ([]*Book, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, title, authors, TO_CHAR(published_date, 'yyyy-mm-dd'),
            page_count, categories, version, created_at, updated_at
		FROM books
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
        AND (categories @> $2 or $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	args := []interface{}{title, pq.Array(categories), filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	books := []*Book{}
	totalRecords := 0

	for rows.Next() {
		var book Book

		err := rows.Scan(
			&totalRecords,
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
			return nil, Metadata{}, err
		}

		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return books, metadata, nil
}
