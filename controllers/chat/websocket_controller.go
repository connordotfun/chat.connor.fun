package chat

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"github.com/aaronaaeng/chat.connor.fun/config"
	_"github.com/aaronaaeng/chat.connor.fun/db/rooms"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/aaronaaeng/chat.connor.fun/context"
	"github.com/satori/go.uuid"
	"github.com/aaronaaeng/chat.connor.fun/db"
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

func makeResponseHeader(ac context.AuthorizedContext) http.Header {
	jwtStr := ac.JWTString()
	if jwtStr != "" {
		return http.Header{
			"Sec-WebSocket-Protocol": []string{jwtStr},
		}
	}
	return nil
}

func HandleWebsocket(hubs *HubMap, roomsRepo db.RoomsRepository, messagesRepo db.MessagesRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !config.Debug && !isOriginValid(c.Request().Header.Get("Origin"), c.Request().Host) {
			return c.NoContent(http.StatusForbidden)
		}
		ac := c.(context.AuthorizedContext)

		roomName := c.Param("room")
		hub, err := lookupHub(roomName, hubs, roomsRepo)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		if hub == nil { //no error, but not found
			return c.NoContent(http.StatusNotFound)
		}

		conn, err := upgrader.Upgrade(c.Response(), c.Request(), makeResponseHeader(c.(context.AuthorizedContext)))
		if err != nil {
			log.Println(err)
			if conn != nil {
				conn.Close()
			}
			return err //upgrade failed
		}


		user := model.User{Id: ac.Requestor().Id, Username: ac.Requestor().Username}
		client := &Client{hub: hub, user: user, canWrite: ac.AccessCode().CanCreate(),
			conn: conn, send: make(chan *model.Message), messagesRepo: messagesRepo}
		client.hub.register <- client

		go client.writer()
		go client.reader()

		return nil
	}
}

func lookupHub(name string, hubs *HubMap, roomsRepo db.RoomsRepository) (*Hub, error) {
	hub, ok := hubs.Load(name)
	if !ok { //hub not in memory, check the database
		room, err := roomsRepo.GetByName(name)
		if err != nil {
			return nil, err
		}
		if room == nil { //create a new room
			room = &model.ChatRoom{Id: uuid.NewV4(), Name: name}
			if err := roomsRepo.Add(room); err != nil {
				return nil, err
			}
		}

		hub = NewHub(room) //init a new hub to activate the room
		go hub.runRoom()
		hubs.Store(name, hub)
	}
	return hub, nil
}
