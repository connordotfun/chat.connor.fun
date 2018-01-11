package controllers

import (
	"github.com/aaronaaeng/chat.connor.fun/db"
	"github.com/labstack/echo"
	"github.com/aaronaaeng/chat.connor.fun/context"
	"net/http"
	"time"
	"github.com/aaronaaeng/chat.connor.fun/model/vericode"
	"github.com/aaronaaeng/chat.connor.fun/model"
)

func getVerificationCode(c echo.Context) string {
	postData := struct {
		Code string
	}{}
	c.Bind(&postData)
	return postData.Code
}

func VerifyUserAccount(repository db.TransactionalRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		verificationsRepo := repository.Verifications()
		rolesRepo := repository.Roles()
		ac := c.(context.AuthorizedContext)

		code := getVerificationCode(c)

		verification, err := verificationsRepo.GetByCode(code)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}

		if verification.UserId != ac.Requestor().Id {
			return c.NoContent(http.StatusForbidden)
		}

		if verification.ExpDate < time.Now().Unix() {
			return c.NoContent(http.StatusGone)
		}

		if !verification.Valid {
			return c.NoContent(http.StatusGone)
		}

		if verification.Purpose != vericode.CodeTypeAccountVerification {
			return c.NoContent(http.StatusBadRequest)
		}

		err = verificationsRepo.Invalidate(code) //this things should really all be done in a transaction
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		err = rolesRepo.RemoveUserRole(verification.UserId, model.RoleUnverified)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		err = rolesRepo.Add(verification.UserId, model.RoleNormal)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusOK)
	}
}
