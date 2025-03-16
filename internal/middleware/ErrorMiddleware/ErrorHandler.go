package errormiddleware

import (
	"SeaChat/pkg/exception"
	"errors"

	"github.com/gofiber/fiber/v2"
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
