package chat

import (
	"github.com/labstack/echo"
	"net/http"

	"github.com/gorilla/websocket"
	"log"
	"github.com/aaronaaeng/chat.connor.fun/config"
	"github.com/slimsag/godocmd/testdata"
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

func HandleWebsocket(c echo.Context) {
	if !isOriginValid(c.Request().Header.Get("Origin"), c.Request().Host) {
		c.NoContent(http.StatusForbidden)
	}
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		c.NoContent(http.StatusBadRequest)
	}
	//TODO: Need a way to store the different rooms without race conditions or bottlenecking
}
