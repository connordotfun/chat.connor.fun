package controllers

import (
	"github.com/labstack/echo"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/aaronaaeng/chat.connor.fun/config"
	"github.com/aaronaaeng/chat.connor.fun/db"
	"github.com/satori/go.uuid"
	"github.com/aaronaaeng/chat.connor.fun/email"
	"github.com/aaronaaeng/chat.connor.fun/model/vericode"
	"strings"
	"errors"
)

func doUserInitNoEmail(u *model.User, repo db.TransactionalRepository) (err error) {
	tx := repo.CreateTransaction()
	defer func() { //if there's an error, rollback
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if err := tx.Users().Add(u); err != nil {
		return err
	}

	if err := tx.Roles().Add(u.Id, model.RoleAnon); err != nil {
		return err
	}

	if err := tx.Roles().Add(u.Id, model.RoleNormal); err != nil {
		return err
	}
	u.Roles = []model.Role{model.Roles.GetRole(model.RoleNormal), model.Roles.GetRole(model.RoleAnon)}

	return nil
}

func doUserInitWithEmail(u *model.User, repo db.TransactionalRepository, host string) (err error) {
	tx := repo.CreateTransaction()
	defer func() { //if there's an error, rollback
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if valid, reason := u.IsEmailValid(); !valid {
		return errors.New("Bad Email: " + reason)
	}

	if err := tx.Users().Add(u); err != nil {
		return err
	}
	if err := tx.Roles().Add(u.Id, model.RoleAnon); err != nil {
		return err
	}

	verification, err := model.GenerateVerificationCode(u.Id, vericode.CodeTypeAccountVerification)
	if err != nil {
		return err
	}

	if err := tx.Verifications().Add(verification); err != nil {
		return err
	}

	if err := tx.Roles().Add(u.Id, model.RoleUnverified); err != nil {
		return err
	}

	u.Roles = []model.Role{model.Roles.GetRole(model.RoleUnverified), model.Roles.GetRole(model.RoleAnon)}

	if err := email.SendAccountVerificationEmail(u.Email, u.Username, makeAccountVerificationLink(host, verification.Code)); err != nil {
		return err
	}

	return nil
}

func CreateUser(repository db.TransactionalRepository, useEmailVerification bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		var u model.User
		if err := c.Bind(&u); err != nil {
			return c.JSON(http.StatusBadRequest, model.Response{
				Error: &model.ResponseError{Type: "BAD_BINDING", Message: err.Error()},
				Data: nil,
			})
		}

		if valid, reason := u.IsUsernameValid(); !valid {
			return c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_USERNAME: " + reason))
		}
		if valid, reason := u.IsSecretValid(); !valid {
			return c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_SECRET: " + reason))
		}

		if foundUser, _ := repository.Users().GetByUsername(u.Username); foundUser != nil {
			return c.JSON(http.StatusBadRequest, model.NewErrorResponse("USER_EXISTS"))
		}
		if foundUser, _ :=repository.Users().GetByEmail(u.Email); foundUser != nil {
			return c.JSON(http.StatusBadRequest, model.NewErrorResponse("EMAIL_IN_USE"))
		}

		u.Id = uuid.NewV4()
		hashedSecret, err := bcrypt.GenerateFromPassword([]byte(u.Secret), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.Response{
				Error: &model.ResponseError{Type: "BAD_PASSWORD", Message: err.Error()},
				Data: nil,
			})
		}
		u.Secret = string(hashedSecret)


		if useEmailVerification {
			err = doUserInitWithEmail(&u, repository, c.Request().Host)
		} else {
			err = doUserInitNoEmail(&u, repository)
		}
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.NewErrorResponse("USER_INIT_FAILED"))
		}

		return c.JSON(http.StatusCreated, model.Response{
			Error: nil,
			Data: model.User{Id: u.Id, Email: u.Email, Username: u.Username, Roles: u.Roles},
		})
	}
}

func GetUser(repository db.TransactionalRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRepo := repository.Users()
		rolesRepo := repository.Roles()
		idStr := c.Param("id")
		id, err := uuid.FromString(idStr)
		if err != nil {
			return c.JSON(http.StatusNotFound, model.NewErrorResponse("BAD_ID"))
		}

		user, err := userRepo.GetById(id)
		if err != nil || user == nil {
			return c.JSON(http.StatusNotFound, model.NewErrorResponse("USER_NOT_FOUND"))
		}

		roles, err := rolesRepo.GetUserRoles(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("RETRIEVE_FAILED"))
		}

		return c.JSON(http.StatusOK, model.NewDataResponse(model.User{Id: user.Id, Username: user.Username, Roles: roles}))
	}
}

func LoginUser(repository db.TransactionalRepository) echo.HandlerFunc {
	return func(c echo.Context) error {  //TODO: generate JWTs
		userRepo := repository.Users()
		rolesRepo := repository.Roles()
		var toLoginUser model.User
		if err := c.Bind(&toLoginUser); err != nil {
			return c.JSON(http.StatusBadRequest, model.Response{
				Error: &model.ResponseError{Type: "BAD_BINDING", Message: err.Error()},
				Data: nil,
			})
		}
		matchedUser, err := userRepo.GetByUsername(toLoginUser.Username)
		if err != nil { //TODO: Handle users not found case
			return c.JSON(http.StatusBadRequest, model.Response{
				Error: &model.ResponseError{Type: "USER_NOT_FOUND", Message: err.Error()},
				Data: nil,
			})
		}

		if matchedUser == nil {
			return c.JSON(http.StatusNotFound, model.NewErrorResponse("USER_NOT_FOUND"))
		}

		if bcrypt.CompareHashAndPassword([]byte(matchedUser.Secret), []byte(toLoginUser.Secret)) != nil {
			return c.JSON(http.StatusUnauthorized, model.Response{
				Error: &model.ResponseError{Type: "PASSWORD_MATCH_FAILED", Message: "Passwords don't match!"},
				Data: nil,
			})
		}
		roles, err := rolesRepo.GetUserRoles(matchedUser.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("ROLE_RETRIEVE_FAILED"))
		}
		userToReturn := model.User{
			Id: matchedUser.Id,
			Email: matchedUser.Email,
			Username: matchedUser.Username,
			Roles: roles,
		}
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

func UpdateUser(repository db.TransactionalRepository) echo.HandlerFunc {
	return nil
}

func makeAccountVerificationLink(host string, code string) string {
	return host + "/verify/account/" + code
}