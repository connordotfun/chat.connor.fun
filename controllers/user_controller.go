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

func handleNewUserInit(u *model.User, usersRepo db.UserRepository, verificationsRepo db.VerificationCodeRepository,
	rolesRepo db.RolesRepository, host string, useEmailVerification bool) error {
	if err := rolesRepo.Add(u.Id, model.RoleUnverified); err != nil {
		return err
	}

	if err := rolesRepo.Add(u.Id, model.RoleAnon); err != nil {
		return err
	}
	if useEmailVerification {
		u.Valid = false
		if !strings.Contains(u.Email,"@") {
			return errors.New("email not valid")
		}
		verification, err := model.GenerateVerificationCode(u.Id, vericode.CodeTypeAccountVerification)
		if err != nil {
			return err
		}
		err = verificationsRepo.Add(verification)
		if err != nil {
			return err
		}
		err = email.SendAccountVerificationEmail(u.Email, u.Username, makeAccountVerificationLink(host, verification.Code))
		if err != nil {
			return err
		}
		if err := rolesRepo.Add(u.Id, model.RoleUnverified); err != nil {
			return err
		}
	} else {
		u.Valid = true
		if err := rolesRepo.Add(u.Id, model.RoleNormal); err != nil {
			return err
		}
	}

	return nil
}

func CreateUser(userRepo db.UserRepository, rolesRepo db.RolesRepository,
		verificationsRepo db.VerificationCodeRepository, useEmailVerification bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		var u model.User
		if err := c.Bind(&u); err != nil {
			return c.JSON(http.StatusBadRequest, model.Response{
				Error: &model.ResponseError{Type: "BAD_BINDING", Message: err.Error()},
				Data: nil,
			})
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

		err = handleNewUserInit(&u, userRepo, verificationsRepo, rolesRepo, c.Request().Host, useEmailVerification)
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.NewErrorResponse("USER_INIT_FAILED"))
		}

		return c.JSON(http.StatusCreated, model.Response{
			Error: nil,
			Data: model.User{Id: u.Id, Email: u.Email, Username: u.Username, Valid: u.Valid},
		})
	}
}

func GetUser(userRepo db.UserRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := uuid.FromString(idStr)
		if err != nil {
			return c.JSON(http.StatusNotFound, model.NewErrorResponse("BAD_ID"))
		}

		user, err := userRepo.GetById(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, model.NewErrorResponse("USER_NOT_FOUND"))
		}

		return c.JSON(http.StatusOK, model.NewDataResponse(model.User{Id: user.Id, Username: user.Username}))
	}
}

func LoginUser(userRepo db.UserRepository) echo.HandlerFunc {
	return func(c echo.Context) error {  //TODO: generate JWTs
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

		if bcrypt.CompareHashAndPassword([]byte(matchedUser.Secret), []byte(toLoginUser.Secret)) != nil {
			return c.JSON(http.StatusUnauthorized, model.Response{
				Error: &model.ResponseError{Type: "PASSWORD_MATCH_FAILED", Message: "Passwords don't match!"},
				Data: nil,
			})
		} else {
			userToReturn := model.User{
				Id: matchedUser.Id,
				Email: matchedUser.Email,
				Username: matchedUser.Username,
				Valid: matchedUser.Valid,
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
}

func UpdateUser(userRepo db.UserRepository) echo.HandlerFunc {
	return nil
}

func makeAccountVerificationLink(host string, code string) string {
	return host + "/verify/account/" + code
}