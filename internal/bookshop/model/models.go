package model

import (
	"database/sql"
	"errors"
)

var (
	ErrEditConflict   = errors.New("edit conflict")
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateEmail = errors.New("duplicate email")
)

type Models struct {
	Books  BookStore
	Users  UserStore
	Tokens TokenStore
}

func NewModels(db *sql.DB) Models {
	return Models{
		Books:  &BookModel{DB: db},
		Users:  &UserModel{DB: db},
		Tokens: &TokenModel{DB: db},
	}
}
