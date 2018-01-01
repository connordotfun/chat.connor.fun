package model


type ChatMessage struct {
	Creator *User `json:"sender,omitempty"`
	Text string `json:"text"`
	CreateDate int64 `json:"createTime"`
	Room ChatRoom
}


