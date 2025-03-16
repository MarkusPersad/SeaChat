package entity

import "github.com/golang-jwt/jwt/v5"

type SeaClaim struct {
	UserID string `json:"user_id"`
	Admin bool `json:"admin"`
	jwt.RegisteredClaims
}
