package model


type ResponseError struct {
	Type string `json:"type"`
	Message string `json:"message"`
}

type Response struct {
	Error *ResponseError `json:"error"`
	Data interface{} `json:"data"`
}

func NewDataResponse(data interface{}) Response {
	return Response{
		Error: nil,
		Data: data,
	}
}

func NewErrorResponse(errorType string) Response {
	return Response{
		Error: &ResponseError{Type: errorType},
		Data: nil,
	}
}