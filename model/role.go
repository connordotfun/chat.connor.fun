package model

import (
	"encoding/json"
)

const (
	RoleAnon = "anonUser"
	RoleNormal = "normalUser"
	RoleUnverified = "unverifiedUser"
	RoleAdmin = "admin"
	RoleBanned = "banned"
)

type Role struct {
	Name 		string 		  `json:"name"`
	Override	string		  `json:"override"`
	Permissions []Permission `json:"permissions"`
}

type RoleMap struct {
	data map[string]Role
}

var Roles RoleMap

func InitRoleMap(rolesJsonData []byte) error {
	err := json.Unmarshal(rolesJsonData, &Roles.data)
	if err != nil {
		return err
	}
	return err
}

func (m RoleMap) GetRole(name string) Role {
	return m.data[name]
}