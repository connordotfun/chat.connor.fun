package chat

import (
	"sync"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/aaronaaeng/chat.connor.fun/db"
)

type Hub struct {
	clients map[*Client]bool

	broadcast chan *model.Message

	register   chan *Client
	unregister chan *Client

	stop chan bool

	Room *model.ChatRoom
	roomsRepo db.RoomsRepository
}

type HubMap struct {
	data sync.Map
}

func NewHubMap() *HubMap {
	return &HubMap{}
}

func (rm *HubMap) Store(roomName string, hub *Hub) {
	rm.data.Store(roomName, hub)
}

func (rm *HubMap) Load(roomName string) (value *Hub, ok bool) {
	res, ok := rm.data.Load(roomName)
	if ok {
		return res.(*Hub), ok
	}
	return nil, ok
}

func (rm *HubMap) Delete(roomName string) {
	rm.data.Delete(roomName)
}

func NewHub(room *model.ChatRoom) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *model.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		stop:       make(chan bool),
		Room:       room,
	}
}

func (r *Hub) runRoom() {
	//there's a possible race condition where the room enters the stop state and a client gets added at the same time
	for {
		select {
		case stop := <-r.stop:
			if stop {
				return
			}
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
				if len(r.clients) == 0 {
					r.stop <- true //stop self
				}
			}
		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- message:
				default: //failed to send, close the client
					close(client.send)
					delete(r.clients, client)
				}

			}
		}
	}
}
