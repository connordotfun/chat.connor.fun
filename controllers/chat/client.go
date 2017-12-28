package chat

import "golang.org/x/net/websocket"

type Client struct {
	room *Room
	conn *websocket.Conn
	send chan []byte
}
