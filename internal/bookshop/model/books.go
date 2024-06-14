package model

import (
	"database/sql"
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
	Insert(book *Book) error
	Get(id int64) (*Book, error)
	Update(book *Book) error
	Delete(id int64) error
}

func (m BookModel) Insert(book *Book) error {
	query := `
		INSERT INTO books (title, authors, published_date, page_count, categories)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, version, created_at, updated_at`

	input := []any{
		book.Title,
		pq.Array(book.Authors),
		book.PublishedDate,
		book.PageCount,
		pq.Array(book.Categories),
	}

	return m.DB.QueryRow(query, input...).Scan(
		&book.ID, &book.Version, &book.CreatedAt, &book.UpdatedAt)
}

func (m BookModel) Get(id int64) (*Book, error) {
	return nil, nil
}

func (m BookModel) Update(book *Book) error {
	return nil
}

func (m BookModel) Delete(id int64) error {
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
