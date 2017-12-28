package chat


type Room struct {
	clients map[*Client]bool

	broadcast chan []byte

	register chan *Client
	unregister chan *Client

	stop chan bool
}


func NewRoom() *Room {
	return &Room{
		clients: make(map[*Client]bool),
		broadcast: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client),
		stop: make(chan bool),
	}
}

func (r *Room) runRoom() {
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