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
)

type testHanlder struct {
	e *echo.Echo
	handler echo.HandlerFunc
	err error
}

func (t *testHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.err = t.handler(t.e.NewContext(r, w))
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

func (c *testWsClient) read() (model.ChatMessage, error) {
	_, bMess, err := c.conn.ReadMessage()
	if err != nil {
		return model.ChatMessage{}, err
	}
	var messages []*model.ChatMessage
	err = json.Unmarshal(bMess, &messages)
	if err != nil {
		return model.ChatMessage{}, err
	}
	return *messages[0], err
}

func TestHandleWebsocket_UpgradeWS(t *testing.T) {
	e := echo.New()

	hubMap := NewHubMap()

	shouldBeTrue := false
	handler := testHanlder{
		e: e,
		handler: func(c echo.Context) error {
			shouldBeTrue = true
			return HandleWebsocket(hubMap, true, c)
		},
	}

	s := httptest.NewServer(&handler)

	d := wstest.NewDialer(&handler, t.Log)

	conn, resp, err := d.Dial("ws://" + s.Listener.Addr().String() + "/ws", nil)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "upgrade failed")

	err = conn.WriteJSON("this won't do anything")
	assert.NoError(t, err)

	assert.NoError(t, handler.err)
	assert.True(t, shouldBeTrue)
}


func TestHandleWebsocket_MultipleClients(t *testing.T) {
	e := echo.New()

	hubMap := NewHubMap()

	shouldBeTrue := false
	handler := testHanlder{
		e: e,
		handler: func(c echo.Context) error {
			shouldBeTrue = true
			return HandleWebsocket(hubMap, true, c)
		},
	}

	s := httptest.NewServer(&handler)

	dClient1 := wstest.NewDialer(&handler, t.Log)
	dClient2 := wstest.NewDialer(&handler, t.Log)

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
	e := echo.New()

	hubMap := NewHubMap()

	shouldBeTrue := false
	handler := testHanlder{
		e: e,
		handler: func(c echo.Context) error {
			shouldBeTrue = true
			return HandleWebsocket(hubMap, true, c)
		},
	}

	s := httptest.NewServer(&handler)

	d := wstest.NewDialer(&handler, t.Log)

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
	e := echo.New()

	hubMap := NewHubMap()

	shouldBeTrue := false
	handler := testHanlder{
		e: e,
		handler: func(c echo.Context) error {
			shouldBeTrue = true
			return HandleWebsocket(hubMap, false, c)
		},
	}

	s := httptest.NewServer(&handler)

	d := wstest.NewDialer(&handler, t.Log)

	conn, resp, err := d.Dial("ws://" + s.Listener.Addr().String() + "/ws", nil)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "upgrade failed")

	client := testWsClient{conn: conn}
	err = client.write("foobar") //writing illegally should cause me to get kicked
	assert.NoError(t, err)

	_, err = client.read()
	assert.Error(t, err)
}