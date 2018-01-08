package testutil

import (
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/satori/go.uuid"
	"errors"
)

type MockRolesRepository struct {
	Roles map[uuid.UUID]map[string]bool
}

func NewMockRolesRepository() *MockRolesRepository {
	return &MockRolesRepository{map[uuid.UUID]map[string]bool{}}
}

func (r *MockRolesRepository) Add(userId uuid.UUID, role string) error {
	if vals, ok := r.Roles[userId]; ok {
		if vals[role] {
			return errors.New("duplicate role")
		}
		vals[role] = true
	} else {
		r.Roles[userId] = map[string]bool{role: true}
	}
	return nil
}

func (r *MockRolesRepository) GetUserRoles(userId uuid.UUID) ([]*model.Role, error) {
	if vals, ok := r.Roles[userId]; ok {
		roles := make([]*model.Role, 0)
		for r := range vals {
			roleObj := model.Roles.GetRole(r)
			roles = append(roles, &roleObj)
		}
		return roles, nil
	} else {
		return make([]*model.Role, 0), nil
	}
}

func (r *MockRolesRepository) RemoveUserRole(userId uuid.UUID, roleName string) error {
	if vals, ok := r.Roles[userId]; ok {
		if roleVal, ok := vals[roleName]; ok && roleVal {
			delete(vals, roleName)
		}
	}
	return nil
}