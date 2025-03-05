package server

import (
	"SeaChat/middleware/fiberswagger"
	"SeaChat/middleware/fiberzerolog"
	"SeaChat/middleware/recoverery"
	"time"

	"github.com/gofiber/contrib/monitor"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/rs/zerolog/log"
)

func(s *FiberServer) RegisterFiberRoutes(){
	// Apply cors middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET","POST","PUT","DELETE","OPTIONS","PATCH"},
		AllowHeaders: []string{"Accept","Authorization","Content-Type"},
		MaxAge: 300,
	}))
	// Apply zerolog middleware
	s.App.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log.Logger,
	}))
	// Apply swagger middleware
	s.App.Use(fiberswagger.New(fiberswagger.Config{
		Title: "SeaChat API Docs",
		BasePath: "/",
		Path: "docs",
		FilePath: "./docs/swagger.json",
		Next: nil,
	}))
	// Apply recover middleware
	s.App.Use(recoverer.New(recoverer.Config{
		Next: nil,
		EnableStackTrace: true,
		StackTraceHandler: recoverery.StackTraceHandler,
	}))
	// Apply metrics middleware
	s.App.Get("/metrics",monitor.New(monitor.Config{
		Title: "SeaChat Metrics Page",
	}))
	// Apply favicon middleware
	s.App.Use(favicon.New(favicon.Config{
		URL: "/favicon.ico",
		File: "pkg/static/chat.svg",
	}))
	// Apply compression middleware
	s.App.Use(compress.New(compress.Config{
		Next: nil,
		Level: compress.LevelBestSpeed,
	}))
	// Apply etag middleware
	s.App.Use(etag.New(etag.ConfigDefault))
	// Apply rate limiter middleware
	s.App.Use(limiter.New(limiter.Config{
		Max: 20,
		Expiration: 20*time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		Next: nil,
	}))
	api := s.App.Group("/api")
	api.Get("/",s.HelloWorldHandler)
}

// HelloWorldHandler
// @Summary Hello World Handler
// @Description Returns a simple hello world message
// @Tags hello
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Router / [get]
func (s *FiberServer) HelloWorldHandler(c fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}
