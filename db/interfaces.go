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

type RoomRepository interface {
	Add(room *model.ChatRoom) error
	GetByName(name string) (*model.ChatRoom, error)
}

type MessagesRepository interface {
	Add(message *model.Message) error
	GetByUserId(userId uuid.UUID) ([]*model.Message, error)
	GetByRoomName(name string) ([]*model.Message, error)
	GetByUserAndRoom(userId int64, name string) ([]*model.Message, error)
}