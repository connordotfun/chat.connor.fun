package auth

import (
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	User       model.User `json:"users"`
	Permissions []model.Permission `json:"permissions,omitempty"`
	*jwt.StandardClaims
}

func (c Claims) Valid() error { //TODO: validate users is some way
	return c.StandardClaims.Valid()
}