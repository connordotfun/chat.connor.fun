package jwtmiddleware

import (
	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"github.com/aaronaaeng/chat.connor.fun/config"
)




const (
	tokenName = "Bearer"
)

func JwtAuth(appConfig config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenStr, err := getJWT(c)
			if err != nil {
				doAuthorization(next, nil, c)
			}

			token, err := parseJWT(tokenStr, appConfig)
			if err != nil {
				doAuthorization(next, nil, c)
			}

			claims, ok := token.Claims.(Claims);
			if !token.Valid || !ok {
				//Handle JWT error case... this should return
			}

			return doAuthorization(next, &claims, c)
		}
	}
}

func parseJWT(tokenStr string, c config.Config) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("expected JWT to be HMAC sgined")
		}

		return []byte(c.JWTSecretKey), nil
	})
}

func getJWT(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	tokenNameLength := len(tokenName)
	if len(authHeader) > tokenNameLength + 1 && authHeader[:tokenNameLength] == tokenName { //If "Bearer: xxxx"
		return authHeader[tokenNameLength + 1:], nil
	}
	return "", errors.New("no JWT bearer token")
}
