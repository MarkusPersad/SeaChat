package fibertimeout

import (
	"SeaChat/pkg/exception"

	"github.com/gofiber/fiber/v3"
)

func TimeoutHandler(_ fiber.Ctx) error {

	return exception.ErrTimeout
}
