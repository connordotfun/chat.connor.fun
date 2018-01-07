package controllers

import (
	"github.com/aaronaaeng/chat.connor.fun/db"
	"github.com/labstack/echo"
	"github.com/aaronaaeng/chat.connor.fun/context"
	"net/http"
	"time"
	"github.com/aaronaaeng/chat.connor.fun/model/vericode"
)

func getVerificationCode(c echo.Context) string {
	postData := struct {
		code string
	}{}
	c.Bind(&postData)
	return postData.code
}

func VerifyUserAccount(verificationsRepo db.VerificationCodeRepository, usersRepo db.UserRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
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

		verificationsRepo.Invalidate(code)

		return c.NoContent(http.StatusOK)
	}
}