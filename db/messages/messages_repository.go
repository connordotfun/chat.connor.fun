package messages

import (
	"github.com/jmoiron/sqlx"
	"github.com/aaronaaeng/chat.connor.fun/model"
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

func (r repository) Add(m *model.Message) (*model.Message, error) {
	params := map[string]interface{} {
		"userId": m.Creator.Id,
		"roomId": m.Room.Id,
		"text": m.Text,
		"createDate": m.CreateDate,
	}

	_, err := r.db.NamedExec(insertMessageQuery, params)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r repository) GetByUserId(userId int64) ([]*model.Message, error) {
	return nil, nil
}