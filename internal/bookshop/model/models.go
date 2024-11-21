package model

import "database/sql"

type Models struct {
	Books BookStore
}

func NewModels(db *sql.DB) Models {
	return Models{
		Books: &BookModel{DB: db},
	}
}
