package utils

import (
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func ResponseSuccess(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ResponseError(message, error string) Response {
	return Response{
		Success: false,
		Message: message,
		Error:   error,
	}
}

func ValidateStruct(s interface{}) error {
	validate := validator.New()
	return validate.Struct(s)
}
