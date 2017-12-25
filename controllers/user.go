package controllers

import (
	"github.com/labstack/echo"
	"github.com/aaronaaeng/chat.connor.fun/db/user"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"net/http"
)


func CreateUser(c echo.Context) error {
	var u model.User
	if err := c.Bind(&u); err != nil {
		return err
	}
	if err := user.Repo.Create(u); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func LoginUser(c echo.Context) error {
	return nil
}