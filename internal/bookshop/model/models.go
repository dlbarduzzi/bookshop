package model

import (
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Books BookStore
}

func NewModels(db *sql.DB) Models {
	return Models{
		Books: &BookModel{DB: db},
	}
}
