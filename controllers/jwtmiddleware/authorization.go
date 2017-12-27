package jwtmiddleware

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/aaronaaeng/chat.connor.fun/model"
)

func doAuthorization(next echo.HandlerFunc, claims *Claims, c echo.Context) error {
	permissions := make([]model.Permission, 0)

	if claims != nil { //there are authenticated claims
		err := claims.Valid()
		if err != nil {
			c.JSON(http.StatusUnauthorized, invalidTokenResponse)
		}
	}

}
