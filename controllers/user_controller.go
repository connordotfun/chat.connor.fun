package controllers

import (
	"github.com/labstack/echo"
	"github.com/aaronaaeng/chat.connor.fun/db/users"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/aaronaaeng/chat.connor.fun/db/roles"
	"github.com/aaronaaeng/chat.connor.fun/config"
)


func CreateUser(c echo.Context) error {
	var u model.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			Error: &model.ResponseError{Type: "BAD_BINDING", Message: err.Error()},
			Data: nil,
		})
	}
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(u.Secret), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Error: &model.ResponseError{Type: "BAD_PASSWORD", Message: err.Error()},
			Data: nil,
		})
	}
	u.Secret = string(hashedSecret)
	createdUser, err := users.Repo.Create(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Error: &model.ResponseError{Type: "USER_CREATE_FAILED", Message: err.Error()},
			Data: nil,
		})
	}
	if err := roles.Repo.AddRole(createdUser.Id, "normal_user"); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Error: &model.ResponseError{Type: "ROLE_ASSIGN_FAILED", Message: err.Error()},
			Data: nil,
		})
	}
	createdUser.Secret = "" //don't return secret
	return c.JSON(http.StatusCreated, model.Response{
		Error: nil,
		Data: createdUser,
	})
}

func LoginUser(c echo.Context) error {  //TODO: generate JWTs
	var toLoginUser model.User
	if err := c.Bind(&toLoginUser); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			Error: &model.ResponseError{Type: "BAD_BINDING", Message: err.Error()},
			Data: nil,
		})
	}
	matchedUser, err := users.Repo.GetByUsername(toLoginUser.Username)
	if err != nil { //TODO: Handle users not found case
		return c.JSON(http.StatusBadRequest, model.Response{
			Error: &model.ResponseError{Type: "USER_NOT_FOUND", Message: err.Error()},
			Data: nil,
		})
	}

	if bcrypt.CompareHashAndPassword([]byte(matchedUser.Secret), []byte(toLoginUser.Secret)) != nil {
		return c.JSON(http.StatusUnauthorized, model.Response{
			Error: &model.ResponseError{Type: "PASSWORD_MATCH_FAILED", Message: "Passwords don't match!"},
			Data: nil,
		})
	} else {
		userToReturn := model.User{Id: matchedUser.Id, Username: matchedUser.Username, Secret: ""}
		jwtStr, err := generateJWT(userToReturn, []byte(config.JWTSecretKey))

		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.Response{
				Error: &model.ResponseError{Type: "JWT_FAILED", Message: "Server failed to sign JWT"},
				Data: nil,
			})
		}

		returnData := map[string]interface{} {
			"token": jwtStr,
			"user": userToReturn,
		}

		return c.JSON(http.StatusOK, model.Response{
			Error: nil,
			Data: returnData,
		})
	}
}