package users

import (
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
	"github.com/aaronaaeng/chat.connor.fun/db"
)

type pgUsersRepository struct {
	db db.DataSource
}

func New(db *sqlx.DB) (*pgUsersRepository, error) {
	_, err := db.Exec(createIfNotExistsQuery)
	if err != nil {
		return nil, err
	}
	return &pgUsersRepository{db}, err
}

func (r pgUsersRepository) NewFromSource(source db.DataSource) db.UserRepository {
	return &pgUsersRepository{db: source}
}

func (r pgUsersRepository) Add(user *model.User) error {
	_, err := r.db.NamedExec(insertUserQuery, user)
	if err != nil {
		return err
	}
	return err
}

func (r pgUsersRepository) Update(user *model.User) error {
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

func (r pgUsersRepository) GetById(id uuid.UUID) (*model.User, error) {
	var user model.User
	rows, err := r.db.NamedQuery(getUserByIdQuery, map[string]interface{}{"id": id})
	defer func() {
		rows.Close()
	}()

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
	defer func() {
		rows.Close()
	}()

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