package controllers

import (
	"github.com/labstack/echo"
	"github.com/aaronaaeng/chat.connor.fun/db"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/satori/go.uuid"
	"net/http"
	"strconv"
	"github.com/labstack/gommon/log"
)

func getMessagesRoom(c echo.Context, messagesRepo db.MessagesRepository, roomIdStr string, count int) error {
	roomId, err := uuid.FromString(roomIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_QUERY"))
	}
	var messages []*model.Message
	if count > 0 {
		messages, err = messagesRepo.GetTopByRoom(roomId, count)
	} else {
		messages, err = messagesRepo.GetByRoomId(roomId)
	}

	if err != nil {
		log.Printf("Failed to retrieve messages: %v", err)
		return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("RETRIEVE_FAILED"))
	}
	return c.JSON(http.StatusOK, model.NewDataResponse(messages))
}

func getMessagesUser(c echo.Context, messagesRepo db.MessagesRepository, userIdStr string, count int) error {
	userId, err := uuid.FromString(userIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_QUERY"))
	}
	var messages []*model.Message
	if count > 0 {
		messages, err = messagesRepo.GetTopByUserId(userId, count)
	} else {
		messages, err = messagesRepo.GetByUserId(userId)
	}

	if err != nil {
		log.Printf("Failed to retrieve messages: %v", err)
		return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("RETRIEVE_FAILED"))
	}
	return c.JSON(http.StatusOK, model.NewDataResponse(messages))
}

func getMessagesUsersAndRoom(c echo.Context, messagesRepo db.MessagesRepository, roomIdStr string, userIdStr string, count int) error {
	userId, err := uuid.FromString(userIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_QUERY"))
	}
	roomId, err := uuid.FromString(roomIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_QUERY"))
	}
	var messages []*model.Message
	if count > 0 {
		messages, err = messagesRepo.GetTopByUserAndRoom(userId, roomId, count)
	} else {
		messages, err = messagesRepo.GetByUserAndRoom(userId, roomId)
	}

	if err != nil {
		log.Printf("Failed to retrieve messages: %v", err)
		return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("RETRIEVE_FAILED"))
	}
	return c.JSON(http.StatusOK, model.NewDataResponse(messages))
}

func GetMessages(repository db.TransactionalRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		messagesRepo := repository.Messages()
		roomIdStr := c.QueryParam("room_id")
		userIdStr := c.QueryParam("user_id")

		count, err := strconv.Atoi(c.QueryParam("count"))
		if err != nil {
			count = -1
		}

		if roomIdStr != "" && userIdStr != "" {
			return getMessagesUsersAndRoom(c, messagesRepo, roomIdStr, userIdStr, count)
		} else if roomIdStr != "" {
			return getMessagesRoom(c, messagesRepo, roomIdStr, count)
		} else if userIdStr != "" {
			return getMessagesUser(c, messagesRepo, userIdStr, count)
		}
		return c.JSON(http.StatusBadRequest, model.NewErrorResponse("NO_QUERY"))
	}
}

func GetMessage(repository db.TransactionalRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		messagesRepo := repository.Messages()
		messagesId, err := uuid.FromString(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_ID"))
		}

		message, err := messagesRepo.GetById(messagesId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("FAILED_RETRIEVE"))
		}
		if message == nil {
			return c.JSON(http.StatusNotFound, model.NewErrorResponse("RESOURCE_NOT_FOUND"))
		}
		return c.JSON(http.StatusOK, model.NewDataResponse(message))
	}
}

func UpdateMessage(repository db.TransactionalRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		messagesRepo := repository.Messages()
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

