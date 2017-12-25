package model

const actionCreate = 0x000F
const actionUpdate = 0x00F0
const actionRead = 0x0F00
const actionDelete = 0xF000

type AccessCode int

type Permission struct {
	Path string
	code AccessCode
}

func (p Permission) setAction(action AccessCode) {
	p.code |= action
}

func (p Permission) setNotAction(action AccessCode)  {
	p.code &= ^action
}

func (p Permission) canDoAction(action AccessCode) bool {
	return (p.code & action) != 0
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

func (p Permission) SetCreate() {
	p.setAction(actionCreate)
}

func (p Permission) SetUpdate()  {
	p.setAction(actionUpdate)
}

func (p Permission) SetRead()  {
	p.setAction(actionRead)
}

func (p Permission) SetDelete()  {
	p.setAction(actionDelete)
}

func (p Permission) SetNotCreate() {
	p.setNotAction(actionCreate)
}

func (p Permission) SetNotUpdate()  {
	p.setNotAction(actionUpdate)
}

func (p Permission) SetNotRead()  {
	p.setNotAction(actionRead)
}

func (p Permission) SetNotDelete()  {
	p.setNotAction(actionDelete)
}