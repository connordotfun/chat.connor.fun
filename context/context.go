package context


import (
	"github.com/labstack/echo"
	"github.com/aaronaaeng/chat.connor.fun/model"
)

type AuthorizedContext struct {
	echo.Context
	JwtString string
	Requestor *model.User
	Code model.AccessCode
}
