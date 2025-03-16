package server

import (
	"SeaChat/pkg/constants"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/rs/zerolog/log"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	// Apply logging middleware
	s.App.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log.Logger,
	}))

	// Apply rate limiting middleware
	s.App.Use(limiter.New(limiter.Config{
		Max: constants.LIMITER_TIMES,
		Expiration: constants.LIMITER_TIME*time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	//Apply monitor middleware
	s.App.Get("/metrics", monitor.New(monitor.Config{Title: fmt.Sprintf("%s Metrics Page",os.Getenv("APP_NAME"))}))

	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)

}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
