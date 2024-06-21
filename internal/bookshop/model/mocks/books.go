package mocks

import (
	"database/sql"
	"errors"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
)

type BookModel struct {
	DB *sql.DB
}

var book = &model.Book{
	ID:    3,
	Title: "Test Title",
}

func (m BookModel) GetAll(title string, categories []string, filters model.Filters) ([]*model.Book, model.Metadata, error) {
	return nil, model.Metadata{}, nil
}

func (m BookModel) Insert(book *model.Book) error {
	if book.Title == "Force Error" {
		return errors.New("forced error")
	}
	return nil
}

func (m BookModel) Get(id int64) (*model.Book, error) {
	switch id {
	case 1:
		return nil, model.ErrRecordNotFound
	case 2:
		return nil, errors.New("forced error")
	default:
		return book, nil
	}
}

func (m BookModel) Update(book *model.Book) error {
	return nil
}

func (m BookModel) Delete(id int64) error {
	return nil
}
