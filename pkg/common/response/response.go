package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data any `json:"data"`
}

func Success( message string, data any) *Response {
	return &Response{
		Code: fiber.StatusOK,
		Message: message,
		Data: data,
	}
}

func Error(code int, message string) *Response {
	return &Response{
		Code: code,
		Message: message,
		Data: nil,
	}
}
