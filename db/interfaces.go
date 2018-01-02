package db

import (
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/satori/go.uuid"
)

type UserRepository interface {
	Add(user *model.User) error
	Update(user *model.User) error
	GetAll() ([]*model.User, error)
	GetById(id uuid.UUID) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
}

type RolesRepository interface {
	Add(userId uuid.UUID, roleName string) error
	GetUserRoles(userId uuid.UUID) ([]*model.Role, error)
}

type RoomsRepository interface {
	Add(room *model.ChatRoom) error
	GetById(id uuid.UUID) (*model.ChatRoom, error)
	GetByName(name string) (*model.ChatRoom, error)
}

type MessagesRepository interface {
	Add(message *model.Message) error
	GetById(id uuid.UUID) (*model.Message, error)
	GetByUserId(userId uuid.UUID) ([]*model.Message, error)
	GetByRoomId(roomId uuid.UUID) ([]*model.Message, error)
	GetByUserAndRoom(userId uuid.UUID, name uuid.UUID) ([]*model.Message, error)
}