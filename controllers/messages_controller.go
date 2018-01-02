package controllers

import (
	"github.com/labstack/echo"
	"github.com/aaronaaeng/chat.connor.fun/db"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/satori/go.uuid"
	"net/http"
	"strconv"
	"github.com/slimsag/godocmd/testdata"
)

func GetMessages(messagesRepo db.MessagesRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomIdStr := c.Param("room")
		roomId, err := uuid.FromString(roomIdStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_ID"))
		}

		var messages []*model.Message
		if countStr := c.QueryParam("count"); countStr != "" {
			count, err := strconv.Atoi(countStr)
			if err != nil {
				return c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_QUERY"))
			}
			messages, err = messagesRepo.GetTopByRoom(roomId, count)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("FAILED_RETRIEVE"))
			}
		} else {
			messages, err = messagesRepo.GetByRoomId(roomId)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("FAILED_RETRIEVE"))
			}
		}

		return c.JSON(http.StatusOK, model.NewDataResponse(messages))
	}
}

func GetMessage(messagesRepo db.MessagesRepository) echo.HandlerFunc {
	return nil
}

func UpdateMessage(messagesRepo db.MessagesRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		messagesId, err := uuid.FromString(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_ID"))
		}

		putData := map[string]interface{} {}
		if err := c.Bind(&putData); err != nil {
			return c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_CONTENT"))
		}

		updatedMessage, err := messagesRepo.Update(messagesId, putData["text"].(string));

		if err != nil {
			c.JSON(http.StatusBadRequest, model.NewErrorResponse("COULD_NOT_UPDATE"))
		}

		return c.JSON(http.StatusOK, model.NewDataResponse(updatedMessage))
	}
}

