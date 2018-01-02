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
	"github.com/aaronaaeng/chat.connor.fun/controllers/jwtmiddleware"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/aaronaaeng/chat.connor.fun/config"
	"github.com/aaronaaeng/chat.connor.fun/db"
	"github.com/slimsag/godocmd/testdata"
)

const (
	testDbHost = "localhost"
	testDbPort = 5432
	testDbName = "connorfuntest_user_controller"
)

var testDb *sqlx.DB

func TestMain(m *testing.M) {
	//New the connorfuntest db and reconnect to the new db
	db, err := sqlx.Open("postgres", "postgresql://localhost:5432?sslmode=disable")
	if err != nil {
		panic("failed to establish db connection")
	}

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", testDbName))
	if err != nil {
		panic(err) //Something went horribly wrong
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", testDbName))
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
	_, err:= db.Query(fmt.Sprintf(`SELECT
			pg_terminate_backend (pg_stat_activity.pid)
			FROM
			pg_stat_activity
			WHERE
			pg_stat_activity.datname = '%s';`, testDbName))

	if err != nil {
		panic("FAILED TO KILL BG CONNECTIONS")
	}

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE %s;", testDbName))
	if err != nil {
		panic("FAILED TO DROP TEST DB: " + err.Error())
	}
}

func initTables() (db.UserRepository, db.RolesRepository, error){
	usersRepo, err := users.New(testDb) //these must be inited in the right order
	if err != nil {
		return nil, nil, err
	}
	rolesRepo, err := roles.New(testDb)

	return usersRepo, rolesRepo, err
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

	testUserJson2 = `
		{"username": "test2", "secret": "test"}
	`
)

func TestCreateUser(t *testing.T) {
	userRepo, rolesRepo, err := initTables()
	assert.NoError(t, err)
	e := echo.New()
	req := httptest.NewRequest("POST", "/api/v1/user", strings.NewReader(testUserJson1))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	createUserFunc := CreateUser(userRepo, rolesRepo)
	assert.NoError(t, createUserFunc(c))
	assert.Equal(t, http.StatusCreated, rec.Code)

	var responseObj model.Response
	err = json.Unmarshal(rec.Body.Bytes(), &responseObj)
	assert.NoError(t, err)
	assert.Nil(t, responseObj.Error)
	assert.Equal(t, responseObj.Data.(map[string]interface{})["username"], "test")

	//recreate

	req = httptest.NewRequest("POST", "/api/v1/user", strings.NewReader(testUserJson1))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()

	c = e.NewContext(req, rec)

	assert.NoError(t, createUserFunc(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response model.Response
	err = json.Unmarshal([]byte(rec.Body.String()), &response)
	assert.NoError(t, err)
	assert.Equal(t, "USER_CREATE_FAILED", response.Error.Type)


	cleanUpTables(t)
}

func TestLoginUser(t *testing.T) {
	userRepo, rolesRepo, err := initTables()
	assert.NoError(t, err)
	config.JWTSecretKey = "secret"
	e := echo.New()
	req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(testUserJson1))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	createUserFunc := CreateUser(userRepo, rolesRepo)
	assert.NoError(t, createUserFunc(c))


	req = httptest.NewRequest("POST", "/api/v1/login", strings.NewReader(testUserJson1))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()

	c = e.NewContext(req, rec)

	loginUserFunc := LoginUser(userRepo)
	assert.NoError(t, loginUserFunc(c))

	var response model.Response
	err = json.Unmarshal([]byte(rec.Body.String()), &response)
	assert.NoError(t, err)

	resData := response.Data.(map[string]interface{})

	assert.NotEmpty(t, resData["token"])
	_, err = jwtmiddleware.ParseAndValidateJWT(resData["token"].(string), []byte("secret"))
	assert.NoError(t, err)

	cleanUpTables(t)
}

func TestGetUser(t *testing.T) {
	userRepo, _, err := initTables()
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(testUserJson1))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	getUserFunc := GetUser(userRepo)

}