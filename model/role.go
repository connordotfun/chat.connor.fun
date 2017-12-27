package model

import (
	"io/ioutil"
	"encoding/json"
)

type Role struct {
	Name 		string 		  `json:"name"`
	Override	string		  `json:"override"`
	Permissions []*Permission `json:"permissions"`
}


type RoleMap struct {
	data map[string]Role
}

var Roles RoleMap

func InitRoleMap(defFile string) error {
	rolesData, err := ioutil.ReadFile(defFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rolesData, &Roles.data)
	if err != nil {
		return err
	}
	return err
}

func (m RoleMap) GetRole(name string) Role {
	return m.data[name]
}