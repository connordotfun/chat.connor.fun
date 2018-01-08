package model

import (
	"github.com/satori/go.uuid"
)

type User struct {
	Id       uuid.UUID  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email 	 string `json:"email,omitempty"`
	Secret   string `json:"secret,omitempty"`
	Roles 	[]Role `json:"roles,omitempty"`
}