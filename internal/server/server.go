package server

import (
	"SeaChat/middleware/fibererrorhandler"
	"os"

	"github.com/gofiber/fiber/v3"
)

type FiberServer struct {
	*fiber.App
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			AppName: os.Getenv("APP_NAME"),
			ServerHeader: os.Getenv("APP_NAME"),
			ErrorHandler: fibererrorhandler.ErrorHandler,
		}),
	}
	return server
}
