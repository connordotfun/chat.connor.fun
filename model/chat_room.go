package model

import "github.com/satori/go.uuid"

type ChatRoom struct {
	Id uuid.UUID `json:"id"`
	Name string `json:"name"`
	Members []User `json:"members"`
	GeoArea	`json:"GeoArea"`
}
