package server

import (
	"github.com/gofiber/fiber/v2"

	"SeaChat/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "SeaChat",
			AppName:      "SeaChat",
		}),

		db: database.New(),
	}

	return server
}
