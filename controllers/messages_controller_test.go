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
)

func TestGetMessage(t *testing.T) {
	messagesRepo := testutil.NewMockMessagesRepository()



	e := echo.New()
	req := httptest.NewRequest("GET", "/api/v1/messages/slakjdakdj", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)


	messageToGet := model.Message{Id: uuid.NewV4()}
	messagesRepo.Add(&messageToGet)

	c.SetParamNames("id")
	c.SetParamValues(messageToGet.Id.String())

	getMessageFunc := GetMessage(messagesRepo)

	assert.NoError(t, getMessageFunc(c))


}
