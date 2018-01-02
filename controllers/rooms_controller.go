package controllers

import (
	"github.com/aaronaaeng/chat.connor.fun/db"
	"github.com/labstack/echo"
	"github.com/aaronaaeng/chat.connor.fun/controllers/chat"
)

func GetNearbyRooms(roomsRepo db.RoomsRepository) echo.HandlerFunc {
	return nil
}

func GetRoomMembers(hubMap *chat.HubMap) echo.HandlerFunc {
	return nil
}
