package model

import (
	"encoding/json"
	"strings"
)

const actionCreate = 0x000F
const actionUpdate = 0x00F0
const actionRead = 0x0F00
const actionDelete = 0xF000

type AccessCode int

func (ac AccessCode) canDoAction(action AccessCode) bool {
	return (ac & action) != 0
}

func (ac AccessCode) CanCreate() bool {
	return ac.canDoAction(actionCreate)
}

func (ac AccessCode) CanRead() bool {
	return ac.canDoAction(actionRead)
}

func (ac AccessCode) CanUpdate() bool {
	return ac.canDoAction(actionUpdate)
}

func (ac AccessCode) CanDelete() bool {
	return ac.canDoAction(actionDelete)
}

type Permission struct {
	Path string
	Code AccessCode
}

func (p *Permission) setAction(action AccessCode) {
	p.Code |= action
}

func (p *Permission) setNotAction(action AccessCode)  {
	p.Code &= ^action
}

func (p Permission) canDoAction(action AccessCode) bool {
	return p.Code.canDoAction(action)
}

func (p Permission) CanCreate() bool {
	return p.canDoAction(actionCreate)
}

func (p Permission) CanUpdate() bool {
	return p.canDoAction(actionUpdate)
}

func (p Permission) CanRead() bool {
	return p.canDoAction(actionRead)
}

func (p Permission) CanDelete() bool {
	return p.canDoAction(actionDelete)
}

func (p *Permission) SetCreate() {
	p.setAction(actionCreate)
}

func (p *Permission) SetUpdate()  {
	p.setAction(actionUpdate)
}

func (p *Permission) SetRead()  {
	p.setAction(actionRead)
}

func (p *Permission) SetDelete()  {
	p.setAction(actionDelete)
}

func (p *Permission) SetNotCreate() {
	p.setNotAction(actionCreate)
}

func (p *Permission) SetNotUpdate()  {
	p.setNotAction(actionUpdate)
}

func (p *Permission) SetNotRead()  {
	p.setNotAction(actionRead)
}

func (p *Permission) SetNotDelete()  {
	p.setNotAction(actionDelete)
}

func (p Permission) generateVerbsStr() string {
	verbs := ""
	if p.CanCreate() {
		verbs += "c"
	}
	if p.CanRead() {
		verbs += "r"
	}
	if p.CanUpdate() {
		verbs += "u"
	}
	if p.CanDelete() {
		verbs += "d"
	}
	return verbs
}

func GenerateVerbCode(verbs string) AccessCode {
	var code AccessCode
	if strings.Contains(verbs, "c") {
		code |= actionCreate
	}
	if strings.Contains(verbs, "r") {
		code |= actionRead
	}
	if strings.Contains(verbs, "u") {
		code |= actionUpdate
	}
	if strings.Contains(verbs, "d") {
		code |= actionDelete
	}
	return code
}

func (p Permission) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Path string `json:"path"`
		Verbs string `json:"verbs"`
	}{
		Path: p.Path,
		Verbs: p.generateVerbsStr(),
	})
}

func (p *Permission) UnmarshalJSON(b []byte) error {
	var perm map[string]string
	err := json.Unmarshal(b, &perm)
	if err != nil {
		return err
	}
	p.Path = perm["path"]
	p.Code = GenerateVerbCode(perm["verbs"])
	return nil
}

func (p Permission) DoesMethodMatch(path string) bool {
	return path == "POST" && p.CanCreate() ||
		   path == "PUT" && p.CanUpdate() ||
		   path == "GET" && p.CanRead() ||
		   path == "DELETE" && p.CanDelete()

}

func (p Permission) DoesPathMatch(path string) bool {
	if path == "/" && path == p.Path {
		return true
	}
	splitPath := strings.Split(path, "/")
	splitPermission := strings.Split(p.Path, "/")
	if len(splitPath) == len(splitPermission) {
		for i := 0; i < len(splitPath); i++ {
			if splitPermission[i] != "*" && splitPath[i] != splitPermission[i] {
				return false
			}
		}
		return true
	}
	return false
}

func (p Permission) IsPermitted(method string, path string) bool {
	pathMatch := p.DoesPathMatch(path)
	methodMatch := p.DoesMethodMatch(method)
	return pathMatch && methodMatch
}

func (p Permission) String() string {
	return "Permission: [" + p.Path + ", " + string(p.Code) + "]"
}