package mocks

import (
	"database/sql"
	"time"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
)

type TokenModel struct {
	DB *sql.DB
}

func (m TokenModel) Save(userID int64, ttl time.Duration, scope string) (*model.Token, error) {
	return nil, nil
}

func (m TokenModel) Insert(token *model.Token) error {
	return nil
}

func (m TokenModel) DeleteAllForUser(userID int64, scope string) error {
	return nil
}
