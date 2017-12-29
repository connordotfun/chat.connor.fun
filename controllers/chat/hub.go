package chat

import (
	"sync"
	"github.com/aaronaaeng/chat.connor.fun/model"
)

type Hub struct {
	clients map[*Client]bool

	broadcast chan []byte

	register chan *Client
	unregister chan *Client

	stop chan bool
}

type HubMap struct {
	data sync.Map
}

func (rm *HubMap) Store(roomName string, hub *Hub) {
	rm.data.Store(roomName, hub)
}

func (rm *HubMap) Load(roomName string) (value *Hub, ok bool) {
	res, ok := rm.data.Load(roomName)
	return res.(*Hub), ok
}


func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		broadcast: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client),
		stop: make(chan bool),
	}
}

func (r *Hub) runRoom() {
	for {
		select {
			case client := <- r.register:
				r.clients[client] = true
			case client := <- r.unregister:
				if _, ok := r.clients[client]; ok {
					delete(r.clients, client)
					close(client.send)
				}
			case message := <- r.broadcast:
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