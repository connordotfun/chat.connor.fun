package user

import (
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/jmoiron/sqlx"
)

var Repo Repository //this must be inited before being used

type Repository struct {
	db *sqlx.DB
}

func Init(database *sqlx.DB) (Repository, error) {
	_, err := database.Exec(CreateIfNotExistsQuery)
	if err != nil {
		return Repository{db:nil}, err
	}
	Repo = Repository{db: database}
	return Repo, nil
}

func (r Repository) Create(user model.User) error {
	_, err := r.db.NamedExec(InsertUserQuery, user)
	return err
}

func (r Repository) Update(user model.User) error {
	return nil
}

func (r Repository) GetAll() ([]*model.User, error) {
	users := make([]*model.User, 0)
	err := r.db.Select(&users, GetAllUsersQuery)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r Repository) GetById(id int64) (*model.User, error) {
	var user model.User
	rows, err := r.db.NamedQuery(GetUserByIdQuery, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}
	err = rows.Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r Repository) Delete(user model.User) error {
	return nil
}