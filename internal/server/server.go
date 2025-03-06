package server

import (
	"SeaChat/internal/database"
	"SeaChat/middleware/fibererrorhandler"
	"os"

	"github.com/gofiber/fiber/v3"
)

type FiberServer struct {
	*fiber.App
	service database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			AppName: os.Getenv("APP_NAME"),
			ServerHeader: os.Getenv("APP_NAME"),
			ErrorHandler: fibererrorhandler.ErrorHandler,
		}),
		service: database.New(),
	}
	return server
}
