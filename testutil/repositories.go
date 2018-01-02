package testutil

import (
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/satori/go.uuid"
	"errors"
)

type MockUserRepository struct {
	Users map[uuid.UUID]model.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{Users: map[uuid.UUID]model.User{}}
}

func (r *MockUserRepository) Add(user *model.User) error {
	for k, v := range r.Users {
		if v.Username == user.Username || k == user.Id {
			return errors.New("duplicate entry")
		}
 	}
 	r.Users[user.Id] = *user
 	return nil
}

func (r *MockUserRepository) Update(user *model.User) error {
	return nil
}

func (r *MockUserRepository) GetAll() ([]*model.User, error) {
	return nil, nil
}

func (r *MockUserRepository) GetById(id uuid.UUID) (*model.User, error) {
	val, ok := r.Users[id]
	if !ok {
		return nil, nil
	}
	toReturn := val
	return &toReturn, nil
}

func (r *MockUserRepository) GetByUsername(username string) (*model.User, error) {
	for _, v := range r.Users {
		if v.Username == username {
			toReturn := v
			return &toReturn, nil
		}
	}
	return nil, nil
}


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