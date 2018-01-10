package controllers

import (
	"testing"
	"github.com/aaronaaeng/chat.connor.fun/testutil"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/satori/go.uuid"
	"github.com/aaronaaeng/chat.connor.fun/model/vericode"
	"github.com/stretchr/testify/assert"
	"github.com/labstack/echo"
	"net/http/httptest"
	"strings"
	"github.com/aaronaaeng/chat.connor.fun/context"
)

func testGenCodeJson(code string) string{
	return `{"code": "` + code + `"}`
}

func TestVerifyUserAccount_SuccessfulVerify(t *testing.T) {
	usersRepo := testutil.NewMockUserRepository()
	rolesRepo := testutil.NewMockRolesRepository()
	verisRepo := testutil.NewMockVerificationsRepo()

	newUser := model.User{
		Id: uuid.NewV4(),
	}

	usersRepo.Add(&newUser)
	rolesRepo.Add(newUser.Id, model.RoleAnon)
	rolesRepo.Add(newUser.Id, model.RoleUnverified)

	verification, err := model.GenerateVerificationCode(newUser.Id, vericode.CodeTypeAccountVerification)
	assert.NoError(t, err)
	assert.NotNil(t, verification)

	verisRepo.Add(verification)

	verifyUserFunc := VerifyUserAccount(verisRepo, usersRepo, rolesRepo)

	e := echo.New()
	req := httptest.NewRequest("POST", "/verifications/accountverification",
		strings.NewReader(testGenCodeJson(verification.Code)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ac := &context.AuthorizedContextImpl{
		Context: e.NewContext(req, rec),
	}
	ac.SetRequestor(newUser)

	assert.NoError(t, verifyUserFunc(ac))
	assert.Equal(t, 200, rec.Code)

	assert.False(t, verisRepo.Data[verification.Code].Valid)
	assert.NotEqual(t, -1, verisRepo.Data[verification.Code].UpdateDate)

	userRoles := rolesRepo.Roles[newUser.Id]
	assert.Equal(t, 2, len(userRoles))
	assert.True(t, userRoles[model.RoleNormal])
	assert.False(t, userRoles[model.RoleUnverified])
}

func TestVerifyUserAccount_WrongUser(t *testing.T) {
	usersRepo := testutil.NewMockUserRepository()
	rolesRepo := testutil.NewMockRolesRepository()
	verisRepo := testutil.NewMockVerificationsRepo()

	newUser := model.User{
		Id: uuid.NewV4(),
	}

	secondUser := model.User{
		Id: uuid.NewV4(),
	}

	usersRepo.Add(&newUser)
	rolesRepo.Add(newUser.Id, model.RoleAnon)
	rolesRepo.Add(newUser.Id, model.RoleUnverified)

	verification, err := model.GenerateVerificationCode(newUser.Id, vericode.CodeTypeAccountVerification)
	assert.NoError(t, err)
	assert.NotNil(t, verification)

	verisRepo.Add(verification)

	verifyUserFunc := VerifyUserAccount(verisRepo, usersRepo, rolesRepo)

	e := echo.New()
	req := httptest.NewRequest("POST", "/verifications/accountverification",
		strings.NewReader(testGenCodeJson(verification.Code)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ac := &context.AuthorizedContextImpl{
		Context: e.NewContext(req, rec),
	}
	ac.SetRequestor(secondUser)

	assert.NoError(t, verifyUserFunc(ac))
	assert.Equal(t, 403, rec.Code)

	userRoles := rolesRepo.Roles[newUser.Id]
	assert.Equal(t, 2, len(userRoles))
	assert.False(t, userRoles[model.RoleNormal])
	assert.True(t, userRoles[model.RoleUnverified])
}

