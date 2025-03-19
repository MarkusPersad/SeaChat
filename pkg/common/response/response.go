package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data any `json:"data"`
	Token string `json:"token"`
}

func Success( message string, data any,token ...string) *Response {
	resp := &Response{
		Code: fiber.StatusOK,
		Message: message,
		Data: data,
	}
	if len(token) >0 {
		resp.Token = token[0]
	} else {
		resp.Token = ""
	}
	return resp
}

func Error(code int, message string) *Response {
	return &Response{
		Code: code,
		Message: message,
		Data: nil,
		Token: "",
	}
}
