package jwtmiddleware

import (
	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"github.com/aaronaaeng/chat.connor.fun/config"
	"net/http"
	"github.com/aaronaaeng/chat.connor.fun/controllers/auth"
	"github.com/aaronaaeng/chat.connor.fun/context"
)

type Skipper func(context echo.Context) bool

type jwtExtractor func(context echo.Context) (string, error)

const (
	tokenName = "Bearer"
)

var (
	defaultExtractor = JWTBearerTokenExtractor
)

func JwtAuth(skipper Skipper, extractor jwtExtractor) echo.MiddlewareFunc {

	if extractor == nil {
		extractor = defaultExtractor
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}
			tokenStr, err := extractor(c)
			if err != nil {
				return doAuthorization(next, nil, c)
			}
			ac := c.(context.AuthorizedContext)
			ac.JwtString = tokenStr

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

func JWTBearerTokenExtractor(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	tokenNameLength := len(tokenName)
	if len(authHeader) > tokenNameLength + 1 && authHeader[:tokenNameLength] == tokenName { //If "Bearer: xxxx"
		return authHeader[tokenNameLength + 1:], nil
	}
	return "", errors.New("no JWT bearer token")
}

func JWTProtocolHeaderExtractor(c echo.Context) (string, error) {
	protocolHeader := c.Request().Header.Get("Sec-WebSocket-Protocol")
	tokenNameLength := len("jwt")
	if len(protocolHeader) > tokenNameLength + 1 && protocolHeader[:tokenNameLength] == "jwt" {
		return protocolHeader[:tokenNameLength + 1], nil
	}
	return "", errors.New("no JWT protocol token")
}
