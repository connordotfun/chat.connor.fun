package chat

import (
	"sync"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/aaronaaeng/chat.connor.fun/db"
	"github.com/labstack/gommon/log"
)

type Hub struct {
	clients map[*Client]bool

	broadcast chan *model.Message

	register   chan *Client
	unregister chan *Client

	stop chan bool

	GetUserList chan (chan []model.User)

	Room *model.ChatRoom
	roomsRepo db.RoomsRepository
	hubMap *HubMap
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

func NewHub(room *model.ChatRoom, hubMap *HubMap) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *model.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		stop:       make(chan bool),
		GetUserList: make(chan chan []model.User),
		Room:       room,
		hubMap: hubMap,
	}
}

func (r *Hub) stopRoom() {
	log.Printf("stopping hub: %s", r.Room.Name)
	for client := range r.clients {
		close(client.send)
		delete(r.clients, client)
	}
	r.hubMap.Delete(r.Room.Name)
}

func (r *Hub) runRoom() {
	defer r.stopRoom()
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
					return
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
		case userListChan := <-r.GetUserList:
			users := make([]model.User, len(r.clients))
			index := 0
			for client := range r.clients {
				users[index] = client.user
			}
			select {
			case userListChan <- users:
			default:
				log.Printf("couldn't send user's list over request")
			}
		}
	}
}
