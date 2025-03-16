package exception

var (
	ErrTimeout  = New(400, "登录超时")
	ErrTokenInvalid = New(401, "Token无效")
)

type SeaError struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

func(err *SeaError) Error() string {
	return err.Message
}

func New(code int, message string) error {
	return &SeaError{
		Code:  code,
		Message: message,
	}
}
