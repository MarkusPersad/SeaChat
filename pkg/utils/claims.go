package utils

import (
	"SeaChat/pkg/constants"
	"SeaChat/pkg/entity"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)
var adminRoles []string

func init(){
	adminString := os.Getenv("ADMIN")
	if strings.TrimSpace(adminString) != "" {
		adminRoles = strings.Split(strings.TrimSpace(adminString), ",")
	}
}

func generateClaims(userid string) entity.SeaClaim {
	return entity.SeaClaim{
		UserID: userid,
		Admin: slices.Contains(adminRoles, userid),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Minute*constants.TOKEN_EXPIRATION),
			},
		},
	}
}

func GenerateTokenString(userid string) (string,error){
	claims := generateClaims(userid)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
