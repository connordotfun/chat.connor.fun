package messages

import (
	"github.com/jmoiron/sqlx"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/satori/go.uuid"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (*repository, error) {
	_, err := db.Exec(createIfNotExistsQuery)
	if err != nil {
		return nil, err
	}
	return &repository{db: db}, nil
}

func (r repository) Add(m *model.Message) error {
	params := map[string]interface{} {
		"id": m.Id,
		"user_id": m.Creator.Id,
		"room_id": m.Room.Id,
		"text": m.Text,
		"create_date": m.CreateDate,
	}

	_, err := r.db.NamedExec(insertMessageQuery, params)
	return err
}

func (r repository) GetById(id uuid.UUID) (*model.Message, error) {
	params := map[string]interface{} {
		"id": id,
	}
	query, err := r.db.PrepareNamed(selectOneByIdQuery)
	if err != nil {
		return nil, err
	}
	message := new(model.Message)
	query.Select(message, params)
	return message, nil
}

func (r repository) GetByUserId(userId uuid.UUID) ([]*model.Message, error) {
	params := map[string]interface{} {
		"user_id": userId,
	}
	query, err := r.db.PrepareNamed(selectByUserIdQuery)
	if err != nil {
		return nil, err
	}
	messages := make([]*model.Message, 0)
	query.Select(&messages, params)

	return messages, err
}

func (r repository) GetByRoomId(roomId uuid.UUID) ([]*model.Message, error) {
	params := map[string]interface{} {
		"room_id": roomId,
	}
	query, err := r.db.PrepareNamed(selectByRoomIdQuery)
	if err != nil {
		return nil, err
	}
	messages := make([]*model.Message, 0)
	query.Select(&messages, params)

	return messages, err
}

func (r repository) GetByUserAndRoom(userId uuid.UUID, roomId uuid.UUID) ([]*model.Message, error) {
	params := map[string]interface{} {
		"user_id": userId,
		"room_id": roomId,
	}
	query, err := r.db.PrepareNamed(selectByUserAndRoomQuery)
	if err != nil {
		return nil, err
	}
	messages := make([]*model.Message, 0)
	query.Select(&messages, params)

	return messages, err
}