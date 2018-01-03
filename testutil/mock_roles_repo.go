package testutil

import (
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/satori/go.uuid"
	"errors"
)

type MockRolesRepository struct {
	Roles map[uuid.UUID][]string
}

func NewMockRolesRepository() *MockRolesRepository {
	return &MockRolesRepository{map[uuid.UUID][]string{}}
}

func (r *MockRolesRepository) Add(userId uuid.UUID, role string) error {
	if vals, ok := r.Roles[userId]; ok {
		for _, r := range vals {
			if r == role {
				return errors.New("duplicate value")
			}
		}
		r.Roles[userId] = append(vals, role)
	} else {
		r.Roles[userId] = []string{role}
	}
	return nil
}

func (r *MockRolesRepository) GetUserRoles(userId uuid.UUID) ([]*model.Role, error) {
	if vals, ok := r.Roles[userId]; ok {
		roles := make([]*model.Role, len(vals))
		for i, r := range vals {
			roleObj := model.Roles.GetRole(r)
			roles[i] = &roleObj
		}
		return roles, nil
	} else {
		return make([]*model.Role, 0), nil
	}
}