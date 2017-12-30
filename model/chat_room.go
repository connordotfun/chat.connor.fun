package model

type ChatRoom struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Members []*User `json:"members"`
	//maybe geolocation data
}
