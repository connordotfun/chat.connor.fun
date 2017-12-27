package roles

import (
	"github.com/jmoiron/sqlx"
	"github.com/aaronaaeng/chat.connor.fun/model"
)

var Repo Repository //this must be inited before being used

type Repository struct {
	db *sqlx.DB
}

func Init(database *sqlx.DB) (Repository, error) {
	_, err := database.Exec(createIfNotExistsQuery)
	if err != nil {
		return Repository{db:nil}, err
	}
	Repo = Repository{db: database}
	return Repo, nil
}


func (r Repository) AddRole(userId int64, role model.Role) error {
	return nil
}

func (r Repository) GetUserRoles(userId int64) error {

	return nil
}