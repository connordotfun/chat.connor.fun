package roles

import (
	"github.com/jmoiron/sqlx"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/satori/go.uuid"
)

type pgRolesRepository struct {
	db *sqlx.DB
}

func New(database *sqlx.DB) (*pgRolesRepository, error) {
	_, err := database.Exec(createIfNotExistsQuery)
	if err != nil {
		return nil, err
	}
	return &pgRolesRepository{db: database}, nil
}


func (r pgRolesRepository) Add(userId uuid.UUID, role string) error {
	params := map[string]interface{} {
		"user_id": userId,
		"role": role,
	}
	_, err := r.db.NamedExec(insertRelationQuery, params)
	return err
}

func (r pgRolesRepository) GetUserRoles(userId uuid.UUID) ([]*model.Role, error) {
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