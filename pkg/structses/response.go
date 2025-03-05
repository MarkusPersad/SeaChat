package structses

import (
	"SeaChat/pkg/exception"

	"github.com/gofiber/fiber/v3"
)

type Response[T any] struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data *T `json:"data"`
}

func CreateResponse[T any](code int, message string, data *T) Response[T] {
	return Response[T]{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func Success[T any](msg string,data *T) Response[T] {
	return Response[T]{
		Code:    fiber.StatusOK,
		Message: msg,
		Data:    data,
	}
}

func Fail[T any](err *exception.PersonalError) Response[T] {
	return Response[T]{
		Code: err.Code,
		Message: err.Error(),
		Data: nil,
	}
}
