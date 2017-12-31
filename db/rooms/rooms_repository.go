package rooms

import (
	"github.com/jmoiron/sqlx"
	"github.com/aaronaaeng/chat.connor.fun/model"
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


func (r pgRoomsRepository) Add(room *model.ChatRoom) (*model.ChatRoom, error){
	_, err := r.db.Exec(insertRoomQuery, &room)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.NamedQuery(selectRoomByNameQuery, &room)
	if err != nil {
		return nil, err
	}

	var insertedRoom model.ChatRoom
	if rows.Next() {
		rows.StructScan(&insertedRoom)
	} else {
		return nil, err //room not found
	}

	return &insertedRoom, nil
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
