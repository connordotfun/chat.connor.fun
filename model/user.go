package model


type User struct {
	Id int64
	Username string
	email string
	secret string
	Roles []*Role
}

