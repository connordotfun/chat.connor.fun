package db

import (
	"github.com/aaronaaeng/chat.connor.fun/model"
	"database/sql"
)


type UserRepository struct {
	db *sql.DB
}

func (r UserRepository) Create(user model.User) error {
	return nil
}

func (r UserRepository) Update(user model.User) error {
	return nil
}

func (r UserRepository) GetAll() ([]*model.User, error) {
	return nil, nil
}

func (r UserRepository) GetById(id int64) (*model.User, error) {
	return nil, nil
}

func (r UserRepository) Delete(user model.User) error {
	return nil
}