package mocks

import (
	"database/sql"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
)

type BookModel struct {
	DB *sql.DB
}

func (m BookModel) GetAll(
	title string,
	categories []string,
	filters model.Filters,
) ([]*model.Book, model.Metadata, error) {
	books := []*model.Book{
		{
			ID:    1,
			Title: "Test Book 1",
		},
	}
	return books, model.Metadata{}, nil
}
