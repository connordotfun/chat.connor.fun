package model

import (
	"github.com/satori/go.uuid"
	"time"
	"crypto/rand" //make sure this is crypto not math
	"encoding/base64"
)

const (
	DefaultLength = 32
	DefaultExpTime = time.Hour * 72

	CodeTypeAccountVerification = "ACCOUNT_VERIFICATION_CODE"
)



type VerificationCode struct {
	Code string `db:"code"`
	Purpose string `db:"purpose"`
	UserId uuid.UUID `db:"user_id"`
	Valid bool `db:"valid"`
	ExpDate int64 `db:"exp_date"`
}

func GenerateVerificationCode(userId uuid.UUID, purpose string) (*VerificationCode, error) {
	expDate := time.Now().Add(time.Duration(DefaultExpTime)).Unix() //valid for 3 days
	return GenerateVerificationCodeWithConfig(userId, purpose, expDate, DefaultLength)
}

func GenerateVerificationCodeWithConfig(userId uuid.UUID, purpose string, expDate int64, length int) (*VerificationCode, error) {

	code, err := generateSecureAscii(length)

	if err != nil {
		return nil, err
	}

	return &VerificationCode{
		Code: code,
		Purpose: purpose,
		UserId: userId,
		Valid: false,
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
