package messages

import (
	"github.com/jmoiron/sqlx"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/satori/go.uuid"
	"github.com/aaronaaeng/chat.connor.fun/db"
)

type repository struct {
	db db.DataSource
}

func New(db *sqlx.DB) (*repository, error) {
	_, err := db.Exec(createIfNotExistsQuery)
	if err != nil {
		return nil, err
	}
	return &repository{db: db}, nil
}

func constructMessageFromJoin(rows *sqlx.Rows) (*model.Message, error) {
	data := &struct {
		Id uuid.UUID
		UserId uuid.UUID `db:"user_id"`
		Username string
		CreateDate int64 `db:"create_date"`
		Text string
		RoomId uuid.UUID `db:"room_id"`
		RoomName string `db:"room_name"`
	}{}

	err := rows.StructScan(&data)
	if err != nil {
		return nil, err
	}

	creator := &model.User{Id: data.UserId, Username: data.Username}
	message := &model.Message{
		Id: data.Id,
		Creator: creator,
		CreateDate: data.CreateDate,
		Text: data.Text,
		Room: &model.ChatRoom{
			Id: data.RoomId,
			Name: data.RoomName,
		},
	}

	return message, nil
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

func (r repository) Update(id uuid.UUID, newText string) (*model.Message, error) {
	params := map[string]interface{} {
		"id": id,
		"text": newText,
	}
	query, err := r.db.PrepareNamed(updateMessageTextQuery)
	if err != nil {
		return nil, err
	}

	resultMessage := new(model.Message)
	err = query.Select(params, resultMessage)
	return resultMessage, err
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

func (r repository) getWithParams(params map[string]interface{}, queryStr string) ([]*model.Message, error) {
	query, err := r.db.PrepareNamed(queryStr)
	if err != nil {
		return nil, err
	}
	messages := make([]*model.Message, 0)
	rows, err := query.Queryx(params)
	for rows.Next() {
		message, err := constructMessageFromJoin(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, err
}

func (r repository) GetByUserId(userId uuid.UUID) ([]*model.Message, error) {
	params := map[string]interface{} {
		"user_id": userId,
	}
	return r.getWithParams(params, selectByUserIdQuery)
}

func (r repository) GetTopByUserId(userId uuid.UUID, count int) ([]*model.Message, error) {
	params := map[string]interface{} {
		"user_id": userId,
		"count": count,
	}
	return r.getWithParams(params, selectTopByUserIdQuery)
}

func (r repository) GetByRoomId(roomId uuid.UUID) ([]*model.Message, error) {
	params := map[string]interface{} {
		"room_id": roomId,
	}
	return r.getWithParams(params, selectByRoomIdQuery)
}

func (r repository) GetTopByRoom(roomId uuid.UUID, count int) ([]*model.Message, error) {
	params := map[string]interface{} {
		"room_id": roomId,
		"count": count,
	}
	return r.getWithParams(params, selectTopByRoomQuery)
}

func (r repository) GetByUserAndRoom(userId uuid.UUID, roomId uuid.UUID) ([]*model.Message, error) {
	params := map[string]interface{} {
		"user_id": userId,
		"room_id": roomId,
	}
	return r.getWithParams(params, selectByUserAndRoomQuery)
}

func (r repository) GetTopByUserAndRoom(userId uuid.UUID, roomId uuid.UUID, count int) ([]*model.Message, error) {
	params := map[string]interface{} {
		"user_id": userId,
		"room_id": roomId,
		"count": count,
	}
	return r.getWithParams(params, selectTopByUserAndRoomQuery)
}