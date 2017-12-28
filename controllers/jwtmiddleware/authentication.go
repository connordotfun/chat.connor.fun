package jwtmiddleware

import (
	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"github.com/aaronaaeng/chat.connor.fun/config"
	"net/http"
	"github.com/aaronaaeng/chat.connor.fun/controllers/auth"
)

type Skipper func(context echo.Context) bool

const (
	tokenName = "Bearer"
)

func JwtAuth(skipper Skipper) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}
			tokenStr, err := getJWT(c)
			if err != nil {
				return doAuthorization(next, nil, c)
			}

			claims, err := ParseAndValidateJWT(tokenStr, []byte(config.JWTSecretKey))
			if err != nil {
				return c.JSON(http.StatusBadRequest, invalidTokenResponse)
			}

			return doAuthorization(next, claims, c)
		}
	}
}

func ParseAndValidateJWT(tokenStr string, jwtSecretKey []byte) (*auth.Claims, error) {

	var claims auth.Claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("expected JWT to be HMAC sgined")
		}

		return jwtSecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("failed to validate token")
	}

	return &claims, nil
}

func getJWT(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	tokenNameLength := len(tokenName)
	if len(authHeader) > tokenNameLength + 1 && authHeader[:tokenNameLength] == tokenName { //If "Bearer: xxxx"
		return authHeader[tokenNameLength + 1:], nil
	}
	return "", errors.New("no JWT bearer token")
}
