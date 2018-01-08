package verifications

import (
	"github.com/jmoiron/sqlx"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"time"
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

func (r *pqVerificationCodeRepository) Invalidate(code string) error {
	params := map[string]interface{} {
		"code": code,
		"update_date": time.Now().Unix(),
	}
	_, err := r.db.NamedExec(invalidateCodeQuery, params)
	return err
}

func (r *pqVerificationCodeRepository) GetByCode(code string) (*model.VerificationCode, error) {
	params := map[string]interface{} {
		"code": code,
	}
	query, err := r.db.PrepareNamed(selectByCodeQuery)
	if err != nil {
		return nil, err
	}
	verificationCode := new(model.VerificationCode)
	rows, err := query.Queryx(params)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		rows.StructScan(verificationCode)
	}
	return verificationCode, err
}