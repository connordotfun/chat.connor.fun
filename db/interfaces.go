package db

import "github.com/aaronaaeng/chat.connor.fun/model"

type UserRepository interface {
	Add(user model.User) (*model.User, error)
	Update(user model.User) error
	GetAll() ([]*model.User, error)
	GetById(id int64) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
}

type RolesRepository interface {
	Add(userId int64, roleName string) error
	GetUserRoles(userId int64) ([]*model.Role, error)
}

type RoomRepository interface {
	Add(room *model.ChatRoom) (*model.ChatRoom, error)
	GetByName(name string) (*model.ChatRoom, error)
}

type MessagesRepository interface {
	Add(message *model.ChatMessage) (*model.ChatMessage, error)
	GetByUserId(userId int64) ([]*model.ChatMessage, error)
	GetByRoomName(name string) ([]*model.ChatMessage, error)
	GetByUserAndRoom(userId int64, name string) ([]*model.ChatMessage, error)
}