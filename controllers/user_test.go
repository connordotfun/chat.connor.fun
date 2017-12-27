package controllers

import (
	"github.com/jmoiron/sqlx"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/aaronaaeng/chat.connor.fun/db/users"
	"testing"
	"github.com/aaronaaeng/chat.connor.fun/db/roles"
	_"github.com/lib/pq"
	"github.com/labstack/echo"
	"net/http/httptest"
	"strings"
	"net/http"
	"encoding/json"
)

const (
	testDbHost = "localhost"
	testDbPort = 5432
	testDbName = "connorfuntest"
)

var testDb *sqlx.DB

func TestMain(m *testing.M) {
	//Init the connorfuntest db and reconnect to the new db
	db, err := sqlx.Open("postgres", "postgresql://localhost:5432?sslmode=disable")
	if err != nil {
		panic("failed to establish db connection")
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS connorfuntest")
	if err != nil {
		panic(err) //Something went horribly wrong
	}

	_, err = db.Exec("CREATE DATABASE connorfuntest;")
	if err != nil {
		panic("Failed to create test db: " + err.Error())
	}

	testDb, err = sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%d dbname=%s sslmode=disable", testDbHost, testDbPort, testDbName))
	if err != nil {
		panic("failed to establish db connection")
	}

	m.Run()

	//clean up the connorfuntestdb
	testDb.Close()
	cleanUpTestDb(db)

	db.Close()
}

func cleanUpTestDb(db *sqlx.DB) {
	_, err:= db.Query(`SELECT
			pg_terminate_backend (pg_stat_activity.pid)
			FROM
			pg_stat_activity
			WHERE
			pg_stat_activity.datname = 'connorfuntest'`)

	if err != nil {
		panic("FAILED TO KILL BG CONNECTIONS")
	}

	_, err = db.Exec("DROP DATABASE connorfuntest;")
	if err != nil {
		panic("FAILED TO DROP TEST DB: " + err.Error())
	}
}

func testRowCountEquals(t *testing.T, expected int) {
	rows, err := testDb.Queryx("SELECT COUNT(*) FROM user_roles")
	assert.NoError(t, err)

	assert.True(t, rows.Next())
	var count int
	err = rows.Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, expected, count)
}

func initTables() {
	_, _ = users.Init(testDb) //these must be inited in the right order
	_, _ = roles.Init(testDb)
}

func cleanUpTables(t *testing.T) {
	_, err := testDb.Exec("DROP TABLE user_roles")
	assert.NoError(t, err)
	_, err = testDb.Exec("DROP TABLE users")
}


const (
	testUserJson1 = `
		{"username": "test", "secret": "test"}
	`
	testUserJsonResponse1 = `{"error":null,"data":{"id":1,"username":"test"}}`
)

func TestCreateUser(t *testing.T) {
	initTables()

	e := echo.New()
	req := httptest.NewRequest("POST", "/api/v1/user", strings.NewReader(testUserJson1))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	assert.NoError(t, CreateUser(c))
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, testUserJsonResponse1, rec.Body.String())

	//recreate

	req = httptest.NewRequest("POST", "/api/v1/user", strings.NewReader(testUserJson1))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()

	c = e.NewContext(req, rec)

	assert.NoError(t, CreateUser(c))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response Response
	err := json.Unmarshal([]byte(rec.Body.String()), &response)
	assert.NoError(t, err)
	assert.Equal(t, "USER_CREATE_FAILED", response.Error.Type)


	cleanUpTables(t)
}