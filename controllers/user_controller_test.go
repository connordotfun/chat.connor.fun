package controllers

import (
	"github.com/stretchr/testify/assert"
	"testing"
	_"github.com/lib/pq"
	"github.com/labstack/echo"
	"net/http/httptest"
	"strings"
	"net/http"
	"encoding/json"
	"github.com/aaronaaeng/chat.connor.fun/controllers/jwtmiddleware"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/aaronaaeng/chat.connor.fun/config"
	"github.com/aaronaaeng/chat.connor.fun/testutil"
	"github.com/satori/go.uuid"
)

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
	userRepo := testutil.NewMockUserRepository()
	rolesRepo := testutil.NewMockRolesRepository()
	verisRepo := testutil.NewMockVerificationsRepo()

	repo := &testutil.MockTransactionalRepository{
		MockRepository: testutil.MockRepository {
			UsersRepo: userRepo,
			RolesRepo: rolesRepo,
			VerificationsRepo: verisRepo,
		},
	}

	e := echo.New()
	req := httptest.NewRequest("POST", "/api/v1/user", strings.NewReader(testUserJson1))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	createUserFunc := CreateUser(repo, false)
	assert.NoError(t, createUserFunc(c))
	assert.Equal(t, http.StatusCreated, rec.Code)

	var responseObj model.Response
	err := json.Unmarshal(rec.Body.Bytes(), &responseObj)
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
	assert.Equal(t, "USER_INIT_FAILED", response.Error.Type)
}

func TestLoginUser(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	rolesRepo := testutil.NewMockRolesRepository()
	verisRepo := testutil.NewMockVerificationsRepo()

	repo := &testutil.MockTransactionalRepository{
		MockRepository: testutil.MockRepository {
			UsersRepo: userRepo,
			RolesRepo: rolesRepo,
			VerificationsRepo: verisRepo,
		},
	}

	config.JWTSecretKey = "secret"
	e := echo.New()
	req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(testUserJson1))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	createUserFunc := CreateUser(repo, false)
	assert.NoError(t, createUserFunc(c))


	req = httptest.NewRequest("POST", "/api/v1/login", strings.NewReader(testUserJson1))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()

	c = e.NewContext(req, rec)

	loginUserFunc := LoginUser(repo)
	assert.NoError(t, loginUserFunc(c))

	var response model.Response
	err := json.Unmarshal([]byte(rec.Body.String()), &response)
	assert.NoError(t, err)

	resData := response.Data.(map[string]interface{})

	assert.NotEmpty(t, resData["token"])
	_, err = jwtmiddleware.ParseAndValidateJWT(resData["token"].(string), []byte("secret"))
	assert.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	rolesRepo := testutil.NewMockRolesRepository()

	repo := &testutil.MockTransactionalRepository{
		MockRepository: testutil.MockRepository {
			UsersRepo: userRepo,
			RolesRepo: rolesRepo,
		},
	}

	userToGet := model.User{Id: uuid.NewV4(), Username: "test", Secret: "Test"}
	userRepo.Add(&userToGet)

	req := httptest.NewRequest("POST", "/api/v1/users/", strings.NewReader(testUserJson1))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userToGet.Id.String())

	getUserFunc := GetUser(repo)

	assert.NoError(t, getUserFunc(c))

	var response model.Response
	err := json.Unmarshal([]byte(rec.Body.String()), &response)
	assert.NoError(t, err)

	userResponse := response.Data.(map[string]interface{})

	userId, err := uuid.FromString(userResponse["id"].(string))
	assert.NoError(t, err)

	assert.Equal(t, userToGet.Id, userId)
	assert.Equal(t, userToGet.Username, userResponse["username"])
}

func TestLoginUser_UserDNE(t *testing.T) {
	usersRepo := testutil.NewMockUserRepository()
	rolesRepo := testutil.NewMockRolesRepository()

	repo := &testutil.MockTransactionalRepository{
		MockRepository: testutil.MockRepository {
			UsersRepo: usersRepo,
			RolesRepo: rolesRepo,
		},
	}

	e := echo.New()
	req := httptest.NewRequest("POST", "/api/v1/login", strings.NewReader(testUserJson1))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	loginUserFunc := LoginUser(repo)
	assert.NoError(t, loginUserFunc(c))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}