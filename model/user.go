package model

import (
	"github.com/satori/go.uuid"
)

type User struct {
	Id       uuid.UUID  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Secret   string `json:"secret,omitempty"`
}