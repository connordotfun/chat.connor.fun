package rooms

import (
	"github.com/jmoiron/sqlx"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/satori/go.uuid"
	"github.com/aaronaaeng/chat.connor.fun/db"
)

type pgRoomsRepository struct {
	db db.DataSource
}

func New(db *sqlx.DB) (*pgRoomsRepository, error) {
	_, err := db.Exec(createIfNotExistsRoomsQuery)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(createCubeQuery)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(createEarthDistQuery)
	if err != nil {
		return nil, err
	}
	return &pgRoomsRepository{db: db}, nil
}

func (r pgRoomsRepository) NewFromSource(source db.DataSource) db.RoomsRepository {
	return &pgRoomsRepository{db: source}
}

func (r pgRoomsRepository) Add(room *model.ChatRoom) error {
	_, err := r.db.NamedExec(insertRoomQuery, &room)
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
	params := map[string]interface{} {
		"name": name,
	}
	rows, err := r.db.NamedQuery(selectRoomByNameQuery, params)
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

func constructRelativeRoom(rows *sqlx.Rows) (*model.RelativeRoom, error) {
	rowData := struct {
		model.ChatRoom
		Distance float64
	}{}

	err := rows.StructScan(&rowData)
	if err != nil {
		return nil, err
	}

	return &model.RelativeRoom{
		Room: model.ChatRoom{
			Id: rowData.Id,
			Name: rowData.Name,
			GeoArea: model.GeoArea {
				Longitude: rowData.Longitude,
				Latitude: rowData.Latitude,
				Radius: rowData.Radius,
			},
		},
		Distance: rowData.Distance,
	}, nil
}

func (r pgRoomsRepository) GetWithinArea(area *model.GeoArea) ([]*model.RelativeRoom, error) {
	rows, err := r.db.NamedQuery(selectWithinRadiusQuery, area)
	if err != nil {
		return nil, err
	}

	rooms := make([]*model.RelativeRoom, 0)
	for rows.Next() {
		resultRoom, err := constructRelativeRoom(rows)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, resultRoom)
	}

	return rooms, nil
}
