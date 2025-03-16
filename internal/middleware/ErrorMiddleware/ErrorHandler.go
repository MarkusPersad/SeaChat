package errormiddleware

import (
	"SeaChat/pkg/exception"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func ErrorHandler(ctx *fiber.Ctx,err error) error {
	if err != nil {
		code := fiber.StatusInternalServerError
		msg := fiber.ErrInternalServerError.Error()
		var e *exception.SeaError
		if errors.As(err, &e){
			code = e.Code
			msg = e.Error()
		}
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"code": code,
			"message": msg,
		})
	}

	return nil
}

func JwtErrorHandler(c *fiber.Ctx, err error) error {
	switch err {
	case jwt.ErrTokenExpired:
		return c.Status(fiber.StatusOK).JSON(exception.ErrTimeout)
	case jwt.ErrTokenInvalidAudience, jwt.ErrTokenInvalidIssuer, jwt.ErrTokenInvalidSubject, jwt.ErrTokenInvalidId, jwt.ErrTokenInvalidClaims:
		return c.Status(fiber.StatusOK).JSON(exception.ErrTokenInvalid)
	default:
		return c.Status(fiber.StatusOK).JSON(fiber.ErrInternalServerError)
	}
}
