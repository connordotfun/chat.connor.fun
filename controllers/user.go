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
		return c.JSON(http.StatusInternalServerError, nil)
	}
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(u.Secret), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	u.Secret = string(hashedSecret)
	userId, err := user.Repo.Create(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	createdUser := model.User{Id: userId, Username: u.Username, Secret: ""}
	return c.JSON(http.StatusCreated, createdUser)
}

func LoginUser(c echo.Context) error {  //TODO: generate JWTs
	var toLoginUser model.User
	if err := c.Bind(&toLoginUser); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	matchedUser, err := user.Repo.GetByUsername(toLoginUser.Username)
	if err != nil { //TODO: Handle user not found case
		return err
	}

	if bcrypt.CompareHashAndPassword([]byte(matchedUser.Secret), []byte(toLoginUser.Secret)) != nil {
		return c.NoContent(http.StatusUnauthorized)
	} else {
		return c.NoContent(http.StatusOK)
	}
}