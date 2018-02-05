package model

import (
	"github.com/satori/go.uuid"
	"strings"
)

const (
	minUsernameLength = 4
	minPasswordLength = 6

)

type User struct {
	Id       uuid.UUID  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email 	 string `json:"email,omitempty"`
	Secret   string `json:"secret,omitempty"`
	Roles 	[]Role `json:"roles,omitempty"`
}

func (u *User) IsUsernameValid() (valid bool, reason string) {
	if len(u.Username) < minUsernameLength {
		return false, "too short"
	}
	if strings.Contains(u.Username, "@") {
		return false, "contains illegal characters"
	}
	return true, ""
}

func (u *User) IsEmailValid() (valid bool, reason string) {
	if len(u.Email) == 0 {
		return false, "email empty"
	}
	if !strings.Contains(u.Email, "@") || !strings.Contains(u.Email, ".") {
		return false, "invalid format"
	}
	return true, ""
}

func (u *User) IsSecretValid() (valid bool, reason string) {
	if len(u.Secret) < minPasswordLength {
		return false, "too short"
	}
	return true, ""
}