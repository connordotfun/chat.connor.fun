package verifications

import (
	"github.com/jmoiron/sqlx"
	"github.com/aaronaaeng/chat.connor.fun/model"
)

type pqVerificationCodeRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (*pqVerificationCodeRepository, error) {
	_, err := db.Exec(createIfNotExistsQuery)
	if err != nil {
		return nil, err
	}
	return &pqVerificationCodeRepository{db}, err
}


func (r *pqVerificationCodeRepository) Add(code *model.VerificationCode) error {
	_, err := r.db.NamedExec(insertCodeQuery, code)
	return err
}