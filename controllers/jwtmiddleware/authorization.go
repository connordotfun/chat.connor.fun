package jwtmiddleware

import (
	"github.com/labstack/echo"
	"github.com/aaronaaeng/chat.connor.fun/controllers"
	"net/http"
)

func doAuthorization(next echo.HandlerFunc, claims *Claims, c echo.Context) error {
	if claims == nil { //No claims = anon users

	}
	err := claims.Valid()
	if err != nil {
		c.JSON(http.StatusUnauthorized, controllers.Response{
			Error: &controllers.ResponseError{Type: "INVALID_AUTH_TOKEN", Message: err.Error()},
			Data: nil,
		})
	}
}
