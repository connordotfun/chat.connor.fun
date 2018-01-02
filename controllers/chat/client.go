package chat

import (
	"github.com/gorilla/websocket"
	"time"
	"github.com/labstack/gommon/log"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"encoding/json"
	_"errors"
	"github.com/satori/go.uuid"
)


const (
	writeWait = 5 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 512
)


type Client struct {
	canWrite bool
	user *model.User
	hub *Hub
	conn *websocket.Conn
	send chan *model.Message
}

func (c *Client) signMessage(messageBytes []byte) (*model.Message, error) {
	var message model.Message
	err := json.Unmarshal(messageBytes, &message)
	if err != nil {
		return nil, err
	}
	message.Id = uuid.NewV4()
	message.CreateDate = time.Now().Unix()
	if c.user != nil {
		message.Creator = &model.User{Id: c.user.Id, Username: c.user.Username}
	} else {
		//return nil, errors.New("message has no creator")
	}
	return &message, nil
}

func (c *Client) writer() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil
	})
	log.Printf("starting client writer")
	for {
		_, message, err := c.conn.ReadMessage()
		log.Printf("received message from client (canWrite: %v)", c.canWrite)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			return
		}
		if c.canWrite {
			signedMessage, err := c.signMessage(message)
			if err != nil {
				log.Printf("error signing message: %v", err)
				return //client sent invalid message, kick
			} else {
				c.hub.broadcast <- signedMessage //eventually we need to process the message
			}
		} else {
			return //client isn't allowed to send, kick
		}
	}
}

func (c *Client) createMessageList(initialMessage *model.Message) ([]byte, error) {
	messageCount := len(c.send) + 1
	messageList := make([]*model.Message, messageCount)

	messageList[0] = initialMessage
	for i := 1; i < messageCount; i++ {
		messageList[i] = <-c.send
	}

	return json.Marshal(messageList)
}

func (c *Client) reader() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok { //send channel is closed
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			toTransmit, err := c.createMessageList(message)
			if err != nil {
				log.Printf("Failed to marshal messages: %v", err)
				return
			}

			w.Write(toTransmit)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

