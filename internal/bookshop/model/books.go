package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dlbarduzzi/bookshop/internal/validator"
	"github.com/lib/pq"
)

type Book struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Authors       []string  `json:"authors"`
	PublishedDate string    `json:"published_date"`
	PageCount     PageCount `json:"page_count"`
	Categories    []string  `json:"categories"`
	Version       int32     `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type BookModel struct {
	DB *sql.DB
}

type BookStore interface {
	GetAll(title string, categories []string, filters Filters) ([]*Book, Metadata, error)
	Insert(book *Book) error
	Get(id int64) (*Book, error)
	Update(book *Book) error
	Delete(id int64) error
}

func (m BookModel) GetAll(title string, categories []string, filters Filters) ([]*Book, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, title, authors, TO_CHAR(published_date, 'yyyy-mm-dd'), page_count,
			categories, version, created_at, updated_at
		FROM books
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
        AND (categories @> $2 or $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	args := []any{title, pq.Array(categories), filters.limit(), filters.offset()}

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

func (m BookModel) Insert(book *Book) error {
	query := `
		INSERT INTO books (title, authors, published_date, page_count, categories)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, version, created_at, updated_at`

	args := []any{
		book.Title,
		pq.Array(book.Authors),
		book.PublishedDate,
		book.PageCount,
		pq.Array(book.Categories),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&book.ID, &book.Version, &book.CreatedAt, &book.UpdatedAt)
}

func (m BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, title, authors, TO_CHAR(published_date, 'yyyy-mm-dd'), page_count,
			categories, version, created_at, updated_at
		FROM books
		WHERE id = $1`

	var book Book

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
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
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &book, nil
}

func (m BookModel) Update(book *Book) error {
	query := `
		UPDATE books
		SET title = $1, authors = $2, published_date = $3, page_count = $4,
			categories = $5, version = version + 1, updated_at = $6
		WHERE id = $7 AND version = $8
		RETURNING version, updated_at`

	args := []any{
		book.Title,
		pq.Array(book.Authors),
		book.PublishedDate,
		book.PageCount,
		pq.Array(book.Categories),
		time.Now().UTC(),
		book.ID,
		book.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&book.Version, &book.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m BookModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM books where id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (b *Book) Validate(v *validator.Validator) {
	b.validateTitle(v)
	b.validateAuthors(v)
	b.validatePublishedDate(v)
	b.validatePageCount(v)
	b.validateCategories(v)
}

func (b *Book) validateTitle(v *validator.Validator) {
	title := strings.TrimSpace(b.Title)

	if title == "" {
		v.AddError("title", "Title is required.")
		return
	}

	if len(title) < 2 || len(title) > 300 {
		v.AddError("title", "Title must have a value between 2 and 300 characters.")
		return
	}

	// Re-assign clean value.
	b.Title = title
}

func (b *Book) validateAuthors(v *validator.Validator) {
	authors := []string{}

	for _, author := range b.Authors {
		author = strings.TrimSpace(author)
		if len(author) < 2 || len(author) > 200 {
			v.AddError("authors", "Authors names must have a value between 2 and 200 characters.")
			return
		} else {
			authors = append(authors, author)
		}
	}

	if len(authors) < 1 {
		v.AddError("authors", "At least 1 author is required.")
		return
	}

	if len(authors) > 3 {
		v.AddError("authors", "Authors cannot have more than 3 values.")
		return
	}

	if !validator.ValuesAreUnique(authors) {
		v.AddError("authors", "Authors cannot have duplicated values.")
		return
	}

	// Re-assign clean value.
	b.Authors = authors
}

func (b *Book) validatePublishedDate(v *validator.Validator) {
	pb := strings.TrimSpace(b.PublishedDate)

	if pb == "" {
		v.AddError("published_date", "Published date is required.")
		return
	}

	t, err := time.Parse("2006-01-02", pb)
	if err != nil {
		v.AddError("published_date", "Invalid published date format (i.e. 2021-02-14).")
		return
	}

	if t.Unix() > time.Now().Unix() {
		v.AddError("published_date", "Published date cannot be in the future.")
		return
	}

	// Re-assign clean value.
	b.PublishedDate = pb
}

func (b *Book) validatePageCount(v *validator.Validator) {
	if b.PageCount < 1 {
		v.AddError("page_count", "Page count must have a value higher than 0.")
		return
	}
}

func (b *Book) validateCategories(v *validator.Validator) {
	categories := []string{}

	for _, category := range b.Categories {
		category = strings.TrimSpace(category)
		if len(category) < 2 || len(category) > 200 {
			v.AddError("categories", "Categories must have values between 2 and 200 characters.")
			return
		} else {
			categories = append(categories, category)
		}
	}

	if len(categories) < 1 {
		v.AddError("categories", "At least 1 category is required.")
		return
	}

	if len(categories) > 5 {
		v.AddError("categories", "Categories cannot have more than 5 values.")
		return
	}

	if !validator.ValuesAreUnique(categories) {
		v.AddError("categories", "Categories cannot have duplicated values.")
		return
	}

	// Re-assign clean value.
	b.Categories = categories
}
