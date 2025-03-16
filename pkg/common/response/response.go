package response

type Response struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data any `json:"data"`
}

func Success(code int, message string, data any) *Response {
	return &Response{
		Code: code,
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
