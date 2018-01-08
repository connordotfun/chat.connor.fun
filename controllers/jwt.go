package controllers

import (
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/aaronaaeng/chat.connor.fun/controllers/auth"
)

func generateJWT(user model.User, jwtSecretKey []byte) (string, error){
	claimedUser := user
	claimedUser.Roles = make([]model.Role, 0)
	claims := auth.Claims{
		User: claimedUser,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Issuer: "connor.fun-login-service",
		},
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	tokenStr, err := token.SignedString(jwtSecretKey) //TODO: this is terrible
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
