package chat

import "github.com/gorilla/websocket"

type Client struct {
	hub *Hub
	conn *websocket.Conn
	send chan []byte
}


func (c *Client) writer() {

}

func (c *Client) reader() {

}

