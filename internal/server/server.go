package server

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"SeaChat/internal/handler"
	errormiddleware "SeaChat/internal/middleware/ErrorMiddleware"
	"SeaChat/internal/model"
)

var (
	AppName = os.Getenv("APP_NAME")
	prefork = os.Getenv("PREFORK") == "true"
)

type FiberServer struct {
	*fiber.App
	*handler.Handler
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: AppName,
			AppName:      AppName,
			Prefork: prefork,
			ErrorHandler: errormiddleware.ErrorHandler,
			ReadTimeout:  30000,
			WriteTimeout: 30000,
			IdleTimeout:  30000,
		}),
		Handler: handler.New(),
	}
	if err := server.InitDB(&model.User{}); err != nil {
		log.Logger.Fatal().Err(err).Msgf("Failed to initialize database Error: %v", err)
	}

	return server
}
