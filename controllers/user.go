package controllers

import (
	"github.com/labstack/echo"
	"github.com/aaronaaeng/chat.connor.fun/db/user"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)


func CreateUser(c echo.Context) error {
	var u model.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Error: ResponseError{Type: "BAD_BINDING", Message: err.Error()},
			Data: nil,
		})
	}
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(u.Secret), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Error: ResponseError{Type: "BAD_PASSWORD", Message: err.Error()},
			Data: nil,
		})
	}
	u.Secret = string(hashedSecret)
	createdUser, err := user.Repo.Create(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Error: ResponseError{Type: "USER_CREATE_FAILED", Message: err.Error()},
			Data: nil,
		})
	}
	createdUser.Secret = "" //don't return secret
	return c.JSON(http.StatusCreated, Response{
		Error: nil,
		Data: createdUser,
	})
}

func LoginUser(c echo.Context) error {  //TODO: generate JWTs
	var toLoginUser model.User
	if err := c.Bind(&toLoginUser); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Error: ResponseError{Type: "BAD_BINDING", Message: err.Error()},
			Data: nil,
		})
	}
	matchedUser, err := user.Repo.GetByUsername(toLoginUser.Username)
	if err != nil { //TODO: Handle user not found case
		return c.JSON(http.StatusBadRequest, Response{
			Error: ResponseError{Type: "USER_NOT_FOUND", Message: err.Error()},
			Data: nil,
		})
	}

	if bcrypt.CompareHashAndPassword([]byte(matchedUser.Secret), []byte(toLoginUser.Secret)) != nil {
		return c.JSON(http.StatusUnauthorized, Response{
			Error: ResponseError{Type: "PASSWORD_MATCH_FAILED", Message: err.Error()},
			Data: nil,
		})
	} else {
		return c.JSON(http.StatusOK, Response{
			Error: nil,
			Data: model.User{Id: matchedUser.Id, Username: matchedUser.Username, Secret: ""},
		})
	}
}