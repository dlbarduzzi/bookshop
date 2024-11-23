package mocks

import (
	"database/sql"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
)

type Models struct {
	Books *BookModel
}

func NewModels(db *sql.DB) model.Models {
	return model.Models{
		Books: &BookModel{DB: db},
	}
}
