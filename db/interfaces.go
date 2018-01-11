package db

import (
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/satori/go.uuid"
	"github.com/jmoiron/sqlx"
	"database/sql"
)


type Repository interface {
	Messages() MessagesRepository
	Users() UserRepository
	Roles() RolesRepository
	Rooms() RoomsRepository
	Verifications() VerificationCodeRepository
}

type DataSource interface {
	sqlx.Queryer
	sqlx.Execer
	sqlx.Preparer
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	Select(dest interface{}, query string, args ...interface{}) error
}

type UserRepository interface {
	Add(user *model.User) error
	Update(user *model.User) error
	GetAll() ([]*model.User, error)
	GetById(id uuid.UUID) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
}

type RolesRepository interface {
	Add(userId uuid.UUID, roleName string) error
	GetUserRoles(userId uuid.UUID) ([]model.Role, error)
	RemoveUserRole(userId uuid.UUID, roleName string) error
}

type RoomsRepository interface {
	Add(room *model.ChatRoom) error
	GetById(id uuid.UUID) (*model.ChatRoom, error)
	GetByName(name string) (*model.ChatRoom, error)
	GetWithinArea(area *model.GeoArea) ([]*model.RelativeRoom, error)
}

type MessagesRepository interface {
	Add(message *model.Message) error
	Update(id uuid.UUID, newText string) (*model.Message, error)
	GetById(id uuid.UUID) (*model.Message, error)

	GetByUserId(userId uuid.UUID) ([]*model.Message, error)
	GetTopByUserId(userId uuid.UUID, count int) ([]*model.Message, error)

	GetByRoomId(roomId uuid.UUID) ([]*model.Message, error)
	GetTopByRoom(roomId uuid.UUID, count int) ([]*model.Message, error)

	GetByUserAndRoom(userId uuid.UUID, name uuid.UUID) ([]*model.Message, error)
	GetTopByUserAndRoom(userId uuid.UUID, name uuid.UUID, count int) ([]*model.Message, error)
}

type VerificationCodeRepository interface {
	Add(code *model.VerificationCode) error
	Invalidate(code string) error
	GetByCode(code string) (*model.VerificationCode, error)
}