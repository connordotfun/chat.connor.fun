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


func (r Repository) AddRole(userId int64, role string) error {
	params := map[string]interface{} {
		"user_id": userId,
		"role": role,
	}
	_, err := r.db.NamedExec(insertRelationQuery, params)
	return err
}

func (r Repository) GetUserRoles(userId int64) ([]*model.Role, error) {
	rows, err := r.db.NamedQuery(getRolesByUserQuery, map[string]interface{}{"user_id": userId})
	if err != nil {
		return nil, err
	}

	userRoles := make([]*model.Role, 0)
	for rows.Next() {
		var roleName string
		err := rows.Scan(&roleName)
		if err != nil {
			return nil, err
		}
		role := model.Roles.GetRole(roleName)
		userRoles = append(userRoles, &role)
	}

	return userRoles, nil
}