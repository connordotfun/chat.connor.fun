package controllers

import (
	"testing"
	"github.com/aaronaaeng/chat.connor.fun/testutil"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/satori/go.uuid"
	"github.com/labstack/echo"
	"net/http/httptest"
	"strings"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"net/http"
)

func TestGetMessage(t *testing.T) {
	repo := testutil.NewEmptyMockTransactionalRepo()

	e := echo.New()
	req := httptest.NewRequest("GET", "/api/v1/messages/slakjdakdj", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)


	messageToGet := model.Message{Id: uuid.NewV4()}
	repo.Messages().Add(&messageToGet)

	c.SetParamNames("id")
	c.SetParamValues(messageToGet.Id.String())

	getMessageFunc := GetMessage(repo)

	assert.NoError(t, getMessageFunc(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	var resMap map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resMap))

	getId, err := uuid.FromString(resMap["data"].(map[string]interface{})["id"].(string))
	assert.NoError(t, err)
	assert.Equal(t, messageToGet.Id, getId)
}

func TestUpdateMessage(t *testing.T) {
	const updateJson = `
	{
		"text": "updated"
	}
	`

	repo := testutil.NewEmptyMockTransactionalRepo()

	e := echo.New()
	req := httptest.NewRequest("GET", "/api/v1/messages/slakjdakdj", strings.NewReader(updateJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)


	messageToGet := model.Message{Id: uuid.NewV4()}
	repo.Messages().Add(&messageToGet)

	c.SetParamNames("id")
	c.SetParamValues(messageToGet.Id.String())

	updateMessageFunc := UpdateMessage(repo)

	assert.NoError(t, updateMessageFunc(c))

	message, _ := repo.Messages().GetById(messageToGet.Id)

	assert.Equal(t, "updated", message.Text)
}

func TestGetMessages_User(t *testing.T) {
	messagesRepo := testutil.NewMockMessagesRepository()

	e := echo.New()
	req := httptest.NewRequest("GET", "/api/v1/messages/slakjdakdj", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	userId := uuid.NewV4()
	userId2 := uuid.NewV4()

	for i := 0; i < 10; i++ {
		messagesRepo.Add(&model.Message{Id: uuid.NewV4(), Creator: &model.User{Id: userId}})
		messagesRepo.Add(&model.Message{Id: uuid.NewV4(), Creator: &model.User{Id: userId2}})
	}

	assert.NoError(t, getMessagesUser(c, messagesRepo, userId.String(), -1))

	var resMap map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resMap))

	results := resMap["data"].([]interface{})

	assert.Len(t, results, 10)
}

func TestGetMessages_UserCount(t *testing.T) {
	messagesRepo := testutil.NewMockMessagesRepository()

	e := echo.New()
	req := httptest.NewRequest("GET", "/api/v1/messages/slakjdakdj", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	userId := uuid.NewV4()

	for i := 0; i < 10; i++ {
		messagesRepo.Add(&model.Message{Id: uuid.NewV4(), Creator: &model.User{Id: userId}})
	}

	assert.NoError(t, getMessagesUser(c, messagesRepo, userId.String(), 3))

	var resMap map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resMap))

	results := resMap["data"].([]interface{})

	assert.Len(t, results, 3)
}

func TestGetMessages_Room(t *testing.T) {
	messagesRepo := testutil.NewMockMessagesRepository()

	e := echo.New()
	req := httptest.NewRequest("GET", "/api/v1/messages/slakjdakdj", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	userId := uuid.NewV4()
	roomId1 := uuid.NewV4()
	roomId2 := uuid.NewV4()

	for i := 0; i < 10; i++ {
		messagesRepo.Add(&model.Message{Id: uuid.NewV4(), Creator: &model.User{Id: userId}, Room: &model.ChatRoom{Id: roomId1}})
		messagesRepo.Add(&model.Message{Id: uuid.NewV4(), Creator: &model.User{Id: userId}, Room: &model.ChatRoom{Id: roomId2}})
	}

	assert.NoError(t, getMessagesRoom(c, messagesRepo, roomId1.String(), -1))

	var resMap map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resMap))

	results := resMap["data"].([]interface{})

	assert.Len(t, results, 10)
}

func TestGetMessages_RoomCount(t *testing.T) {
	messagesRepo := testutil.NewMockMessagesRepository()

	e := echo.New()
	req := httptest.NewRequest("GET", "/api/v1/messages/slakjdakdj", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	userId := uuid.NewV4()
	roomId1 := uuid.NewV4()
	roomId2 := uuid.NewV4()

	for i := 0; i < 10; i++ {
		messagesRepo.Add(&model.Message{Id: uuid.NewV4(), Creator: &model.User{Id: userId}, Room: &model.ChatRoom{Id: roomId1}})
		messagesRepo.Add(&model.Message{Id: uuid.NewV4(), Creator: &model.User{Id: userId}, Room: &model.ChatRoom{Id: roomId2}})
	}

	assert.NoError(t, getMessagesRoom(c, messagesRepo, roomId1.String(), 4))

	var resMap map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resMap))

	results := resMap["data"].([]interface{})

	assert.Len(t, results, 4)
}

func TestGetMessages_Both(t *testing.T) {
	messagesRepo := testutil.NewMockMessagesRepository()

	e := echo.New()
	req := httptest.NewRequest("GET", "/api/v1/messages/slakjdakdj", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	userId := uuid.NewV4()
	userId2 := uuid.NewV4()
	roomId1 := uuid.NewV4()
	roomId2 := uuid.NewV4()

	for i := 0; i < 10; i++ {
		messagesRepo.Add(&model.Message{Id: uuid.NewV4(), Creator: &model.User{Id: userId}, Room: &model.ChatRoom{Id: roomId1}})
		messagesRepo.Add(&model.Message{Id: uuid.NewV4(), Creator: &model.User{Id: userId2}, Room: &model.ChatRoom{Id: roomId1}})
		messagesRepo.Add(&model.Message{Id: uuid.NewV4(), Creator: &model.User{Id: userId}, Room: &model.ChatRoom{Id: roomId2}})
	}

	assert.NoError(t, getMessagesUsersAndRoom(c, messagesRepo, roomId1.String(), userId.String(), -1))

	var resMap map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resMap))

	results := resMap["data"].([]interface{})

	assert.Len(t, results, 10)
}

func TestGetMessages_BothCount(t *testing.T) {
	messagesRepo := testutil.NewMockMessagesRepository()

	e := echo.New()
	req := httptest.NewRequest("GET", "/api/v1/messages/slakjdakdj", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	userId := uuid.NewV4()
	userId2 := uuid.NewV4()
	roomId1 := uuid.NewV4()
	roomId2 := uuid.NewV4()

	for i := 0; i < 10; i++ {
		messagesRepo.Add(&model.Message{Id: uuid.NewV4(), Creator: &model.User{Id: userId}, Room: &model.ChatRoom{Id: roomId1}})
		messagesRepo.Add(&model.Message{Id: uuid.NewV4(), Creator: &model.User{Id: userId2}, Room: &model.ChatRoom{Id: roomId1}})
		messagesRepo.Add(&model.Message{Id: uuid.NewV4(), Creator: &model.User{Id: userId}, Room: &model.ChatRoom{Id: roomId2}})
	}

	assert.NoError(t, getMessagesUsersAndRoom(c, messagesRepo, roomId1.String(), userId.String(), 8))

	var resMap map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resMap))

	results := resMap["data"].([]interface{})

	assert.Len(t, results, 8)
}