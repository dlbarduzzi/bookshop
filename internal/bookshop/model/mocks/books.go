package mocks

import (
	"database/sql"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
)

type BookModel struct {
	DB *sql.DB
}

func (m BookModel) Insert(book *model.Book) error {
	return nil
}

func (m BookModel) Get(id int64) (*model.Book, error) {
	return nil, nil
}

func (m BookModel) Update(book *model.Book) error {
	return nil
}

func (m BookModel) Delete(id int64) error {
	return nil
}
