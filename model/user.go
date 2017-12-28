package model

type User struct {
	Id       int64  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Secret   string `json:"secret,omitempty"`
	roles    []*Role
}

func (u User) Roles() []*Role {
	return u.roles //eventually go to the db for this
}
