package exception

var (
	ErrTimeout = New(400, "请求超时")
)


type PersonalError struct {
	Code int
	Message string
}

func New(code int, message string) error {
	return &PersonalError{
		Code:    code,
		Message: message,
	}
}

func (e *PersonalError) Error() string {
	return e.Message
}
