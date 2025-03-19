package request

type UserInfo struct {
	Info string `json:"info" validate:"required" field_error_info:"用户信息不能为空"`
}