package controllers

import (
	"github.com/aaronaaeng/chat.connor.fun/db"
	"github.com/labstack/echo"
	"github.com/aaronaaeng/chat.connor.fun/controllers/chat"
	"net/http"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"errors"
	"github.com/satori/go.uuid"
	"github.com/aaronaaeng/chat.connor.fun/util"
)

func GetNearbyRooms(repository db.TransactionalRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomsRepo := repository.Rooms()

		latStr := c.QueryParam("lat")
		lonStr := c.QueryParam("lon")
		radiusStr := c.QueryParam("radius")

		parser := util.FloatParser{}

		lat := parser.ParseFloat(latStr, 64)
		lon := parser.ParseFloat(lonStr, 64)
		radius := parser.ParseFloat(radiusStr, 64)

		if parser.Err != nil {
			return c.JSON(http.StatusBadRequest, "BAD_QUERY")
		}

		searchArea := model.GeoArea{
			Latitude: lat,
			Longitude: lon,
			Radius: radius,
		}

		nearbyRooms, err := roomsRepo.GetWithinArea(&searchArea)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("RETRIEVE_FAILED"))
		}

		return c.JSON(http.StatusOK, nearbyRooms)
	}
}

func getRoomMembersList(hub *chat.Hub) ([]model.User, error) {
	if hub == nil {
		return make([]model.User, 0), nil
	}

	usersChan := make(chan []model.User)

	hub.GetUserList <- usersChan

	users, ok := <- usersChan
	if !ok {
		return nil, errors.New("room in illegal state (shutting-down or shutdown)")
	}
	close(usersChan)

	return users, nil
}

func GetRoomMembers(hubMap *chat.HubMap) echo.HandlerFunc { //there's no good way to test this rn
	return func(c echo.Context) error {
		roomName := c.Param("room")
		hub, ok := hubMap.Load(roomName)
		if !ok {
			return c.JSON(http.StatusOK, model.NewDataResponse(make([]*model.User, 0)))
		}

		users, err := getRoomMembersList(hub)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("HUB_SHUTDOWN"))
		}
		return c.JSON(http.StatusOK, model.NewDataResponse(users))
	}
}

func GetRoom(repository db.TransactionalRepository, hubMap *chat.HubMap) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomsRepository := repository.Rooms()
		roomName := c.Param("room")
		hub, ok := hubMap.Load(roomName)
		var room *model.ChatRoom
		if ok {
			room = hub.Room
		} else {
			roomFromDb, err := roomsRepository.GetByName(roomName)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("RETRIEVE_FAILED"))
			}
			if roomFromDb == nil {
				return c.JSON(http.StatusNotFound, model.NewErrorResponse("ROOM_DNE"))
			}
			room = roomFromDb
		}

		roomMembers, err := getRoomMembersList(hub)
		if err != nil {
			roomMembers = nil
		}

		roomToReturn := *room
		roomToReturn.Members = roomMembers

		return c.JSON(http.StatusOK, model.NewDataResponse(roomToReturn))
	}
}

func CreateRoom(repository db.TransactionalRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		rooms := repository.Rooms()
		var chatRoom model.ChatRoom
		if err := c.Bind(&chatRoom); err != nil {
			return c.JSON(http.StatusBadRequest, model.NewErrorResponse("BIND_FAILED"))
		}

		chatRoom.Id = uuid.NewV4()

		if err := rooms.Add(&chatRoom); err != nil {
			return c.JSON(http.StatusBadRequest, model.NewErrorResponse("ROOM_EXISTS"))
		}

		return c.JSON(http.StatusCreated, model.NewDataResponse(chatRoom))
	}
}