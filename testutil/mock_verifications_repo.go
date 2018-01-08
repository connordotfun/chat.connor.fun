package testutil

import (
	"github.com/aaronaaeng/chat.connor.fun/model"
	"errors"
	"time"
)

type MockVerificationsRepo struct {
	Data map[string]model.VerificationCode
}

func NewMockVerificationsRepo() *MockVerificationsRepo {
	return &MockVerificationsRepo{Data: map[string]model.VerificationCode{}}
}

func (r *MockVerificationsRepo) Add(code *model.VerificationCode) error {
	if _, ok := r.Data[code.Code]; ok {
		return errors.New("duplicate entry")
	}
	r.Data[code.Code] = *code
	return nil
}

func (r *MockVerificationsRepo) Invalidate(code string) error {
	if val, ok := r.Data[code]; ok {
		val.Valid = false
		val.UpdateDate = time.Now().Unix()
		r.Data[code] = val
	}
	return nil
}

func (r *MockVerificationsRepo) GetByCode(code string) (*model.VerificationCode, error) {
	if val, ok := r.Data[code]; ok {
		return &val, nil
	}
	return nil, nil
}
