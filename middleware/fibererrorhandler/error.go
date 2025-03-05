package fibererrorhandler

import (
	"SeaChat/pkg/exception"
	"SeaChat/pkg/structses"

	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(ctx fiber.Ctx,err error) error {
	response := structses.Response[any]{}
	if personalError,ok := err.(*exception.PersonalError); ok {
		response = structses.Fail[any](personalError)
		ctx.Status(fiber.StatusOK).JSON(response)
	}
	response = structses.Response[any]{
		Code: fiber.StatusInternalServerError,
		Message: fiber.ErrInternalServerError.Message,
		Data: nil,
	}
	ctx.Status(fiber.StatusInternalServerError).JSON(response)
	return nil
}
