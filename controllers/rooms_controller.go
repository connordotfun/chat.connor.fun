package controllers

import (
	"github.com/aaronaaeng/chat.connor.fun/db"
	"github.com/labstack/echo"
	"github.com/aaronaaeng/chat.connor.fun/controllers/chat"
	"net/http"
	"github.com/aaronaaeng/chat.connor.fun/model"
)

func GetNearbyRooms(roomsRepo db.RoomsRepository) echo.HandlerFunc {
	return nil
}

func GetRoomMembers(hubMap *chat.HubMap) echo.HandlerFunc { //there's no good way to test this rn
	return func(c echo.Context) error {
		roomName := c.Param("room")
		hub, ok := hubMap.Load(roomName)
		if !ok {
			return c.JSON(http.StatusOK, model.NewDataResponse(make([]*model.User, 0)))
		}

		usersChan := make(chan []model.User)
		hub.GetUserList <- usersChan

		users, ok := <- usersChan
		if !ok {
			return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("HUB_SHUTDOWN"))
		}
		close(usersChan)
		return c.JSON(http.StatusOK, model.NewDataResponse(users))
	}
}
