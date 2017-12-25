package model


type Role struct {
	Name string `json:"name"`
}

const normalUser = "NORMAL_USER"
const admin = "ADMIN"

func NormalUser() Role {
	return Role{Name: normalUser}
}

func Admin() Role {
	return Role{Name: admin}
}