package model

import (
	"strings"
)

const CREATE  = "c"
const UPDATE  = "u"
const READ = "r"
const DELETE  = "d"

type Permission struct {
	Path string
	code string
}


func (p Permission) canCreate() bool {
	return strings.Contains(p.code, CREATE)
}

func (p Permission) canUpdate() bool {
	return strings.Contains(p.code, UPDATE)
}

func (p Permission) canRead() bool {
	return strings.Contains(p.code, READ)
}

func (p Permission) canDelete() bool {
	return strings.Contains(p.code, DELETE)
}