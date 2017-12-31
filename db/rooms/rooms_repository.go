package rooms

import (
	"github.com/jmoiron/sqlx"
	"github.com/aaronaaeng/chat.connor.fun/model"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (*Repository, error) {
	_, err := db.Exec(createIfNotExistsRoomsQuery)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}


func (r Repository) Add(room *model.ChatRoom) (*model.ChatRoom, error){
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

func (r Repository) GetByName(name string) (*model.ChatRoom, error) {
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
