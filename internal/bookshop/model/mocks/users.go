package mocks

import (
	"database/sql"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
)

type UserModel struct {
	DB *sql.DB
}

var users = []model.User{
	{
		ID:    1,
		Name:  "Test User",
		Email: "test.user1@email.com",
	},
	{
		ID:    2,
		Name:  "Test User",
		Email: "test.user2@email.com",
	},
}

func (m UserModel) Insert(user *model.User) (*model.User, error) {
	return &users[0], nil
}

func (m UserModel) Update(user *model.User) error {
	return nil
}

func (m UserModel) GetByEmail(email string) (*model.User, error) {
	if email == "test.user2@email.com" {
		return &users[1], nil
	} else {
		return nil, nil
	}
}

func (m UserModel) GetForToken(token string, scope string) (*model.User, error) {
	return nil, nil
}

func (m UserModel) SavePassword(userID int64, passwordHash []byte) error {
	return nil
}
