package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestUser_IsEmailValid_Pass(t *testing.T) {
	u := User{
		Username: "foobar",
		Email: "validemail@gmail.com",
		Secret: "foobar",
	}

	valid, _ := u.IsEmailValid()
	assert.True(t, valid)
}

func TestUser_IsEmailValid_Fail1(t *testing.T) {
	u := User{
		Username: "foobar",
		Email: "validemailgmail.com",
		Secret: "foobar",
	}

	valid, _ := u.IsEmailValid()
	assert.False(t, valid)
}

func TestUser_IsEmailValid_Fail2(t *testing.T) {
	u := User{
		Username: "foobar",
		Email: "validemail@gmailcom",
		Secret: "foobar",
	}

	valid, _ := u.IsEmailValid()
	assert.False(t, valid)
}

func TestUser_IsEmailValid_Fail_Empty(t *testing.T) {
	u := User{
		Username: "foobar",
		Email: "",
		Secret: "foobar",
	}

	valid, _ := u.IsEmailValid()
	assert.False(t, valid)
}

func TestUser_IsUsernameValid(t *testing.T) {
	u := User{
		Username: "foobar",
		Email: "validemail@gmail.com",
		Secret: "foobar",
	}

	valid, _ := u.IsUsernameValid()
	assert.True(t, valid)
}

func TestUser_IsUsernameValid_Fail_Email(t *testing.T) {
	u := User{
		Username: "validemail@gmail.com",
		Email: "validemail@gmail.com",
		Secret: "foobar",
	}

	valid, _ := u.IsUsernameValid()
	assert.False(t, valid)
}

func TestUser_IsUsernameValid_Fail_TooShort(t *testing.T) {
	u := User{
		Username: "foo",
		Email: "validemail@gmail.com",
		Secret: "foobar",
	}

	valid, _ := u.IsUsernameValid()
	assert.False(t, valid)
}

func TestUser_IsSecretValid(t *testing.T) {
	u := User{
		Username: "foobar",
		Email: "validemail@gmail.com",
		Secret: "foobar123",
	}

	valid, _ := u.IsSecretValid()
	assert.True(t, valid)
}

func TestUser_IsSecretValid_Fail_TooShort(t *testing.T) {
	u := User{
		Username: "foobar",
		Email: "validemail@gmail.com",
		Secret: "foo",
	}

	valid, _ := u.IsSecretValid()
	assert.False(t, valid)
}