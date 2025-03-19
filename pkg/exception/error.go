package exception

var (
	ErrTimeout  = New(400, "登录超时")
	ErrTokenInvalid = New(402, "Token无效")
	ErrBadRequest = New(401, "请求参数错误")
	ErrCaptchaInvalid = New(403, "验证码无效")
	ErrUserAlreadyExists = New(405, "用户已存在")
	ErrUserNotFound = New(406, "用户不存在")
	ErrUserStatusInvalid = New(408, "用户状态无效")
	ErrPasswordInvalid = New(407, "密码错误")
	ErrPermissionDenied = New(409, "权限不足")
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
