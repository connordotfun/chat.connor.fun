package context


import (
	"github.com/labstack/echo"
	"github.com/aaronaaeng/chat.connor.fun/model"
)


type AuthorizedContext interface {
	echo.Context

	GetAccessCode() model.AccessCode
	SetAccessCode(code model.AccessCode)

	GetRequestor() *model.User
	SetRequestor(user *model.User)

	GetJWTString() string
	SetJWTString(string string)
}

type AuthorizedContextImpl struct {
	jwtString string
	requestor *model.User
	code model.AccessCode
	echo.Context
}

func (c *AuthorizedContextImpl) GetAccessCode() model.AccessCode {
	return c.code
}

func (c *AuthorizedContextImpl) GetRequestor() *model.User {
	return c.requestor
}

func (c *AuthorizedContextImpl) GetJWTString() string {
	return c.jwtString
}

func (c *AuthorizedContextImpl) SetAccessCode(code model.AccessCode) {
	c.code = code
}

func (c *AuthorizedContextImpl) SetRequestor(user *model.User) {
	c.requestor = user
}

func (c *AuthorizedContextImpl) SetJWTString(jwtString string) {
	c.jwtString = jwtString
}