package model

import (
	"github.com/satori/go.uuid"
	"time"
	"crypto/rand" //make sure this is crypto not math
	"encoding/base64"
	"github.com/aaronaaeng/chat.connor.fun/model/vericode"
)

const (
	DefaultLength = 32
	DefaultExpTime = time.Hour * 72
)



type VerificationCode struct {
	Code string `db:"code"`
	Purpose vericode.VerificationCodeType `db:"purpose"`
	UserId uuid.UUID `db:"user_id"`
	Valid bool `db:"valid"`
	ExpDate int64 `db:"exp_date"`
}

func GenerateVerificationCode(userId uuid.UUID, purpose vericode.VerificationCodeType) (*VerificationCode, error) {
	expDate := time.Now().Add(time.Duration(DefaultExpTime)).Unix() //valid for 3 days
	return GenerateVerificationCodeWithConfig(userId, purpose, expDate, DefaultLength)
}

func GenerateVerificationCodeWithConfig(userId uuid.UUID, purpose vericode.VerificationCodeType, expDate int64, length int) (*VerificationCode, error) {

	code, err := generateSecureAscii(length)

	if err != nil {
		return nil, err
	}

	return &VerificationCode{
		Code: code,
		Purpose: purpose,
		UserId: userId,
		Valid: true,
		ExpDate: expDate,
	}, nil
}

func generateSecureAscii(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
