package model


type User struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Secret string
	Roles []*Role `json:"roles"`
}

