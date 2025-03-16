package server

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"SeaChat/internal/database"
	errormiddleware "SeaChat/internal/middleware/ErrorMiddleware"
)

var (
	AppName = os.Getenv("APP_NAME")
	prefork = os.Getenv("PREFORK") == "true"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: AppName,
			AppName:      AppName,
			Prefork: prefork,
			ErrorHandler: errormiddleware.ErrorHandler,
		}),

		db: database.New(),
	}

	return server
}
