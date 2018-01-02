package rooms

import (
	"github.com/jmoiron/sqlx"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/satori/go.uuid"
)

type pgRoomsRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (*pgRoomsRepository, error) {
	_, err := db.Exec(createIfNotExistsRoomsQuery)
	if err != nil {
		return nil, err
	}
	return &pgRoomsRepository{db: db}, nil
}


func (r pgRoomsRepository) Add(room *model.ChatRoom) error {
	_, err := r.db.Exec(insertRoomQuery, &room)
	if err != nil {
		return err
	}
	return err
}

func (r pgRoomsRepository) GetById(id uuid.UUID) (*model.ChatRoom, error) {
	params := map[string]interface{} {
		"id": id,
	}
	query, err := r.db.PrepareNamed(selectRoomByIdQuery)
	if err != nil {
		return nil, err
	}
	chatRoom := new(model.ChatRoom)
	query.Select(chatRoom, params)
	return chatRoom, nil
}

func (r pgRoomsRepository) GetByName(name string) (*model.ChatRoom, error) {
	rows, err := r.db.NamedQuery(selectRoomByNameQuery, model.ChatRoom{Name: name})
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		var room model.ChatRoom
		rows.StructScan(&room)
		return &room, nil
	}
	return nil, nil //not found
}
