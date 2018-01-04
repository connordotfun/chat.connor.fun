package controllers

import (
	"testing"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/stretchr/testify/assert"
	"github.com/aaronaaeng/chat.connor.fun/controllers/jwtmiddleware"
	"github.com/satori/go.uuid"
)


var (
	testJwtUser1 = model.User{Id: uuid.NewV4(), Username: "test", Secret: "secret"}
)

func TestGenerateJWT(t *testing.T) {
	jwtStr, err := generateJWT(testJwtUser1, []byte("secret"))
	assert.NoError(t, err)
	assert.NotEmpty(t, jwtStr)

	claims, err := jwtmiddleware.ParseAndValidateJWT(jwtStr, []byte("secret"))
	assert.NoError(t, err)
	assert.NotEmpty(t, claims)
}

func TestGenerateJWT_BadKey(t *testing.T) {
	jwtStr, err := generateJWT(testJwtUser1, []byte("wrong-secret"))
	assert.NoError(t, err)
	assert.NotEmpty(t, jwtStr)

	claims, err := jwtmiddleware.ParseAndValidateJWT(jwtStr, []byte("secret"))
	assert.Error(t, err)
	assert.Empty(t, claims)
}
