package utils

import (
	"SeaChat/internal/database"
	"SeaChat/pkg/constants"
	"SeaChat/pkg/entity"
	"SeaChat/pkg/exception"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
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

func TokenCheck(ctx *fiber.Ctx, service database.ValkeyService, admin bool) (string, error) {
	// 初始化返回值
	var tokenString string

	// 从上下文中获取 JWT token 和 claims
	token := ctx.Locals(constants.JWT_CONTEXT_KEY).(*jwt.Token)
	claims := token.Claims.(*entity.SeaClaim)

	// 验证用户是否存在
	if service.GetValue(ctx.UserContext(), constants.JWT_CONTEXT_KEY+":"+claims.UserID) == "" {
		return "", exception.ErrTokenInvalid
	}

	// 检查 token 是否即将过期
	timeLeft := time.Until(claims.ExpiresAt.Time)
	if timeLeft < constants.TOKEN_CHECK*time.Minute {
		// 生成新的 token
		newToken, err := GenerateTokenString(claims.UserID)
		if err != nil {
			log.Error().Err(err).Msg("Generate token string failed")
			return "", err
		}

		// 如果需要管理员权限但用户不是管理员，返回权限不足错误
		if admin && !claims.Admin {
			return "", exception.ErrPermissionDenied
		}

		// 返回新 token
		tokenString = newToken
	}

	return tokenString, nil
}