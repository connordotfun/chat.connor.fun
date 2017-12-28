package model


type ResponseError struct {
	Type string `json:"type"`
	Message string `json:"message"`
}

type Response struct {
	Error *ResponseError `json:"error"`
	Data interface{} `json:"data"`
}