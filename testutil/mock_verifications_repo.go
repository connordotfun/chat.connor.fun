package testutil

import (
	"github.com/aaronaaeng/chat.connor.fun/model"
	"errors"
	"time"
)

type MockVerificationsRepo struct {
	data map[string]model.VerificationCode
}

func (r *MockVerificationsRepo) Add(code *model.VerificationCode) error {
	if _, ok := r.data[code.Code]; ok {
		return errors.New("duplicate entry")
	}
	r.data[code.Code] = *code
	return nil
}

func (r *MockVerificationsRepo) Invalidate(code string) error {
	if val, ok := r.data[code]; ok {
		val.Valid = false
		val.UpdateDate = time.Now().Unix()
		r.data[code] = val
	}
	return nil
}

func (r *MockVerificationsRepo) GetByCode(code string) (*model.VerificationCode, error) {
	if val, ok := r.data[code]; ok {
		return &val, nil
	}
	return nil, nil
}
