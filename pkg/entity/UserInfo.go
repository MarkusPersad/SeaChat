package entity

type UserDetails struct {
	UserID string `json:"user_id"`
	UserName string `json:"user_name"`
	Email string `json:"email"`
	Avatar string `json:"avatar"`
	Status string `json:"status"`
}