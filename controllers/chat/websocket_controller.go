package chat

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"github.com/aaronaaeng/chat.connor.fun/config"
	"github.com/aaronaaeng/chat.connor.fun/db/rooms"
	"github.com/aaronaaeng/chat.connor.fun/model"
)


var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize: 1024,
}

func isOriginValid(origin string, host string) bool {
	var expected string
	if config.Debug {
		expected = "http://" + host
	} else {
		expected = "https://" + host
	}
	return expected == origin
}

func HandleWebsocket(hubs *HubMap, c echo.Context) error {
	if !isOriginValid(c.Request().Header.Get("Origin"), c.Request().Host) {
		return c.NoContent(http.StatusForbidden)
	}

	roomName := c.Param("room")
	hub, err := lookupHub(roomName, hubs)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	if hub == nil { //no error, but not found
		return c.NoContent(http.StatusNotFound)
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		if conn != nil {
			conn.Close()
		}
		return err //upgrade failed
	}

	client := &Client{hub: hub, conn: conn, send: make(chan *model.ChatMessage)}
	client.hub.register <- client

	go client.writer()
	go client.reader()

	return nil
}

func lookupHub(name string, hubs *HubMap) (*Hub, error) {
	hub, ok := hubs.Load(name)
	if !ok { //hub not in memory, check the database
		room, err := rooms.Repo.GetByName(name)
		if err != nil {
			return nil, err
		}
		if room == nil { //room does not exist
			return nil, nil
		}
		hub = NewHub() //init a new hub to activate the room
		//TODO: run the hub
		hubs.Store(name, hub)
	}
	return hub, nil
}
