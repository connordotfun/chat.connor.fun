package users

import (
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/jmoiron/sqlx"
)

type pgUsersRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (*pgUsersRepository, error) {
	_, err := db.Exec(createIfNotExistsQuery)
	if err != nil {
		return nil, err
	}
	return &pgUsersRepository{db}, err
}

func (r pgUsersRepository) Add(user model.User) (*model.User, error) {
	_, err := r.db.NamedExec(insertUserQuery, user)
	if err != nil {
		return nil, err
	}
	return r.GetByUsername(user.Username)
}

func (r pgUsersRepository) Update(user model.User) error {
	return nil
}

func (r pgUsersRepository) GetAll() ([]*model.User, error) {
	users := make([]*model.User, 0)
	err := r.db.Select(&users, getAllUsersQuery)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r pgUsersRepository) GetById(id int64) (*model.User, error) {
	var user model.User
	rows, err := r.db.NamedQuery(getUserByIdQuery, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, nil //No such user
}

func (r pgUsersRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	rows, err := r.db.NamedQuery(getUserByUsernameQuery, map[string]interface{}{"username": username})
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, nil
}