package model


type User struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
	Secret string
	roles []*Role
}

func (u User) Roles() []*Role {
	return u.roles //eventually go to the db for this
}