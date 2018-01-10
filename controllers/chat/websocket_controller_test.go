package chat

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/posener/wstest"
	"github.com/gorilla/websocket"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"encoding/json"
	"time"
	"github.com/aaronaaeng/chat.connor.fun/context"
	"github.com/aaronaaeng/chat.connor.fun/testutil"
	"github.com/aaronaaeng/chat.connor.fun/filter"
)

type testHandler struct {
	e *echo.Echo
	handler echo.HandlerFunc
	code model.AccessCode
	err error
}

func newTestHandler(e *echo.Echo, handler echo.HandlerFunc, code model.AccessCode) *testHandler {
	return &testHandler{e, handler, code, nil}
}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &context.AuthorizedContextImpl{
		Context: t.e.NewContext(r, w),
	}

	c.SetAccessCode(t.code)
	t.err = t.handler(c)
}

type testWsClient struct {
	conn *websocket.Conn
}

func (c *testWsClient) write(message string) error {
	jsonMessage := []byte(`{"text": "` + message + `"}`)
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonMessage)
	if err != nil {
		return err
	}
	return w.Close()
}

func (c *testWsClient) read() (model.Message, error) {
	_, bMess, err := c.conn.ReadMessage()
	if err != nil {
		return model.Message{}, err
	}
	var messages []*model.Message
	err = json.Unmarshal(bMess, &messages)
	if err != nil {
		return model.Message{}, err
	}
	return *messages[0], err
}

func TestHandleWebsocket_UpgradeWS(t *testing.T) {
	messagesRepo := testutil.NewMockMessagesRepository()
	roomsRepo := testutil.NewMockRoomsRepository()

	e := echo.New()

	hubMap := NewHubMap()
	handleWsFunc := HandleWebsocket(hubMap, roomsRepo, messagesRepo, filter.NewTree("../../assets/bannedList.txt"))

	shouldBeTrue := false
	handlerFunc := func(c echo.Context) error {
		shouldBeTrue = true
		return handleWsFunc(c)
	}
	handler := newTestHandler(e, handlerFunc, model.GenerateVerbCode("cr"))

	s := httptest.NewServer(handler)

	d := wstest.NewDialer(handler, nil)

	conn, resp, err := d.Dial("ws://" + s.Listener.Addr().String() + "/ws", nil)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "upgrade failed")

	err = conn.WriteJSON("this won't do anything")
	assert.NoError(t, err)

	assert.NoError(t, handler.err)
	assert.True(t, shouldBeTrue)
}


func TestHandleWebsocket_MultipleClients(t *testing.T) {
	messagesRepo := testutil.NewMockMessagesRepository()
	roomsRepo := testutil.NewMockRoomsRepository()

	e := echo.New()

	hubMap := NewHubMap()
	handleWsFunc := HandleWebsocket(hubMap, roomsRepo, messagesRepo, filter.NewTree("../../assets/bannedList.txt"))

	shouldBeTrue := false
	handlerFunc := func(c echo.Context) error {
		shouldBeTrue = true
		return handleWsFunc(c)
	}
	handler := newTestHandler(e, handlerFunc, model.GenerateVerbCode("cr"))

	s := httptest.NewServer(handler)

	dClient1 := wstest.NewDialer(handler, nil)
	dClient2 := wstest.NewDialer(handler, nil)

	client1Conn, resp, err := dClient1.Dial("ws://" + s.Listener.Addr().String() + "/ws", nil)
	assert.NoError(t, err)

	assert.NoError(t, handler.err)
	assert.True(t, shouldBeTrue)

	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "upgrade failed")

	client1 := testWsClient{conn: client1Conn}

	client2Conn, resp, err := dClient2.Dial("ws://" + s.Listener.Addr().String() + "/ws", nil)
	assert.NoError(t, err)

	assert.NoError(t, handler.err)
	assert.True(t, shouldBeTrue)

	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "upgrade failed")

	client2 := testWsClient{conn: client2Conn}

	defer func() {
		client1.conn.Close()
		client2.conn.Close()
	}()

	err = client1.write("hello")
	assert.NoError(t, err)

	client2.conn.SetReadDeadline(time.Now().Add(time.Duration(5 * time.Second)))
	message, err := client2.read()
	client2.conn.SetReadDeadline(time.Now().Add(time.Duration(5 * time.Second)))
	assert.NoError(t, err)
	assert.Equal(t, "hello", message.Text)
}

func TestHandleWebsocket_IllegalMessage(t *testing.T) {
	messagesRepo := testutil.NewMockMessagesRepository()
	roomsRepo := testutil.NewMockRoomsRepository()

	e := echo.New()

	hubMap := NewHubMap()
	handleWsFunc := HandleWebsocket(hubMap, roomsRepo, messagesRepo, filter.NewTree("../../assets/bannedList.txt"))

	shouldBeTrue := false
	handlerFunc := func(c echo.Context) error {
		shouldBeTrue = true
		return handleWsFunc(c)
	}
	handler := newTestHandler(e, handlerFunc, model.GenerateVerbCode("cr"))

	s := httptest.NewServer(handler)

	d := wstest.NewDialer(handler, nil)

	conn, resp, err := d.Dial("ws://" + s.Listener.Addr().String() + "/ws", nil)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "upgrade failed")

	err = conn.WriteJSON("this will cause me to get kicked by the server")
	assert.NoError(t, err)

	assert.NoError(t, handler.err)
	assert.True(t, shouldBeTrue)

	_, _, err = conn.ReadMessage()
	assert.Error(t, err)

}

func TestHandleWebsocket_ReadOnly(t *testing.T) {
	messagesRepo := testutil.NewMockMessagesRepository()
	roomsRepo := testutil.NewMockRoomsRepository()

	e := echo.New()

	hubMap := NewHubMap()
	handleWsFunc := HandleWebsocket(hubMap, roomsRepo, messagesRepo, filter.NewTree("../../assets/bannedList.txt"))

	shouldBeTrue := false
	handlerFunc := func(c echo.Context) error {
		shouldBeTrue = true
		return handleWsFunc(c)
	}
	handler := newTestHandler(e, handlerFunc, model.GenerateVerbCode("r"))
	s := httptest.NewServer(handler)

	d := wstest.NewDialer(handler, nil)

	conn, resp, err := d.Dial("ws://" + s.Listener.Addr().String() + "/ws", nil)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "upgrade failed")

	client := testWsClient{conn: conn}
	err = client.write("foobar") //writing illegally should cause me to get kicked
	assert.NoError(t, err)

	_, err = client.read()
	assert.Error(t, err)
}

func TestHandleWebsocket_MessageCleaning(t *testing.T) {
	messagesRepo := testutil.NewMockMessagesRepository()
	roomsRepo := testutil.NewMockRoomsRepository()

	e := echo.New()

	hubMap := NewHubMap()
	handleWsFunc := HandleWebsocket(hubMap, roomsRepo, messagesRepo, filter.NewTree("../../assets/bannedList.txt"))

	shouldBeTrue := false
	handlerFunc := func(c echo.Context) error {
		shouldBeTrue = true
		return handleWsFunc(c)
	}
	handler := newTestHandler(e, handlerFunc, model.GenerateVerbCode("cr"))

	s := httptest.NewServer(handler)

	dClient := wstest.NewDialer(handler, nil)

	clientConn, resp, err := dClient.Dial("ws://" + s.Listener.Addr().String() + "/ws", nil)
	assert.NoError(t, err)

	assert.NoError(t, handler.err)
	assert.True(t, shouldBeTrue)

	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "upgrade failed")

	client := testWsClient{conn: clientConn}

	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "upgrade failed")


	defer func() {
		client.conn.Close()
	}()
	err = client.write("apple")
	assert.NoError(t, err)

	client.conn.SetReadDeadline(time.Now().Add(time.Duration(5 * time.Second)))
	message, err := client.read()
	client.conn.SetReadDeadline(time.Now().Add(time.Duration(5 * time.Second)))
	assert.NoError(t, err)
	assert.Equal(t, "*****", message.Text)
}