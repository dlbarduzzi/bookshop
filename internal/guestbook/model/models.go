package model

import "database/sql"

type Models struct {
	Guests GuestStore
}

func NewModels(db *sql.DB) Models {
	return Models{
		Guests: &GuestModel{DB: db},
	}
}
