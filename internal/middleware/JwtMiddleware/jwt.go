package jwtmiddleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)



func JwtFilter(ctx *fiber.Ctx ) bool {
	return strings.Contains(ctx.Path(), "login") || strings.Contains(ctx.Path(), "register") || strings.Contains(ctx.Path(), "getcaptcha")||
		strings.Contains(ctx.Path(), "docs") || strings.Contains(ctx.Path(), "metrics") || strings.Contains(ctx.Path(), "health")
}
