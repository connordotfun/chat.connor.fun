package jwtmiddleware

import (
	"github.com/aaronaaeng/chat.connor.fun/model"
)

var (
	invalidTokenResponse = model.Response{
		Error: &model.ResponseError{Type: "INVALID_AUTH_TOKEN", Message: "Invalid Token"},
		Data: nil,
	}
)