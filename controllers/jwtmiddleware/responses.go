package jwtmiddleware

import "github.com/aaronaaeng/chat.connor.fun/controllers"

var (
	invalidTokenResponse = controllers.Response{
		Error: &controllers.ResponseError{Type: "INVALID_AUTH_TOKEN", Message: "Invalid Token"},
		Data: nil,
	}
)