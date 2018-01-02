package model

import "github.com/satori/go.uuid"

type Message struct {
	Id uuid.UUID `json:"id"`
	Creator *User `json:"sender,omitempty"`
	Text string `json:"text"`
	CreateDate int64 `json:"createTime"`
	Room *ChatRoom
}

