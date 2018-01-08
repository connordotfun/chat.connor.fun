package testutil

import (
	"github.com/satori/go.uuid"
	"github.com/aaronaaeng/chat.connor.fun/model"
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
