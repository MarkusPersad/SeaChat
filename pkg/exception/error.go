package exception


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
